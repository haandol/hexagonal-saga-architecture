package repository

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/constant/status"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/util"
	"gorm.io/gorm"
)

type TripRepository struct {
	BaseRepository
}

func NewTripRepository(db *gorm.DB) *TripRepository {
	return &TripRepository{
		BaseRepository: BaseRepository{DB: db},
	}
}

func (r *TripRepository) PublishStartSaga(ctx context.Context,
	corrID string, parentID string, d *dto.Trip,
) error {
	db := r.WithContext(ctx)

	evt := &command.StartSaga{
		Message: message.Message{
			Name:          reflect.ValueOf(command.StartSaga{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.StartSagaBody{
			TripID:   d.ID,
			CarID:    d.CarID,
			HotelID:  d.HotelID,
			FlightID: d.FlightID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}

	v, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	row := &entity.Outbox{
		KafkaTopic: "saga-service",
		KafkaKey:   evt.CorrelationID,
		KafkaValue: v,
	}
	return db.Create(row).Error
}

func (r *TripRepository) PublishAbortSaga(ctx context.Context,
	corrID string, parentID string, tripID uint, reason string,
) error {
	db := r.WithContext(ctx)

	evt := &command.AbortSaga{
		Message: message.Message{
			Name:          reflect.ValueOf(command.AbortSaga{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.AbortSagaBody{
			TripID: tripID,
			Reason: reason,
			Source: "trip",
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}

	v, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	row := &entity.Outbox{
		KafkaTopic: "saga-service",
		KafkaKey:   evt.CorrelationID,
		KafkaValue: v,
	}
	return db.Create(row).Error
}

func (r *TripRepository) Create(ctx context.Context, corrID, parentID string, d *dto.Trip) (dto.Trip, error) {
	panicked := true

	tx := r.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return dto.Trip{}, err
	}
	defer func() {
		if r := recover(); r != nil || panicked {
			tx.Rollback()
		}
	}()

	txCtx := context.WithValue(ctx, constant.TX("tx"), tx)

	row := entity.Trip{
		UserID:   d.UserID,
		CarID:    d.CarID,
		HotelID:  d.HotelID,
		FlightID: d.FlightID,
		Status:   status.TripInitialized,
	}
	result := tx.Create(&row)
	if result.Error != nil {
		return dto.Trip{}, result.Error
	}

	trip := row.DTO()
	if err := r.PublishStartSaga(txCtx, corrID, parentID, &trip); err != nil {
		return dto.Trip{}, err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return dto.Trip{}, err
	}

	return row.DTO(), nil
}

func (r *TripRepository) Update(ctx context.Context, d *dto.Trip) error {
	result := r.WithContext(ctx).
		Where("id = ? AND user_id = ?", d.ID, d.UserID).
		Updates(&entity.Trip{
			ID:       d.ID,
			UserID:   d.UserID,
			CarID:    d.CarID,
			HotelID:  d.HotelID,
			FlightID: d.FlightID,
			Status:   d.Status,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowAffected
	}
	return nil
}

func (r *TripRepository) List(ctx context.Context) ([]dto.Trip, error) {
	var rows entity.Trips
	result := r.WithContext(ctx).
		Limit(10).
		Order("id desc").
		Find(&rows)
	if result.Error != nil {
		return []dto.Trip{}, result.Error
	}

	return rows.DTO(), nil
}

func (r *TripRepository) Complete(ctx context.Context, evt *event.SagaEnded) error {
	return r.WithContext(ctx).
		Where("id = ?", evt.Body.TripID).
		Updates(&entity.Trip{
			CarBookingID:    evt.Body.CarBookingID,
			HotelBookingID:  evt.Body.HotelBookingID,
			FlightBookingID: evt.Body.FlightBookingID,
			Status:          status.TripBooked,
		}).Error
}

func (r *TripRepository) Abort(ctx context.Context, evt *event.SagaAborted) error {
	return r.WithContext(ctx).
		Where("id = ?", evt.Body.TripID).
		Updates(&entity.Trip{
			Status: status.TripCanceled,
		}).Error
}

func (r *TripRepository) GetByID(ctx context.Context, id uint) (dto.Trip, error) {
	row := &entity.Trip{}
	result := r.WithContext(ctx).
		Where("id = ?", id).
		Take(row)
	if result.Error != nil {
		return dto.Trip{}, result.Error
	}

	return row.DTO(), nil
}
