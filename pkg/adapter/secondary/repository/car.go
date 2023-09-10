package repository

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	ErrNoCarBookingFound = errors.New("no car-booking found")
)

type CarRepository struct {
	BaseRepository
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		BaseRepository{DB: db},
	}
}

func (r *CarRepository) PublishCarBooked(ctx context.Context,
	corrID string, parentID string, d *dto.CarBooking,
) error {
	db := r.WithContext(ctx)

	evt := &event.CarBooked{
		Message: message.Message{
			Name:          reflect.ValueOf(event.CarBooked{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarBookedBody{
			BookingID: d.ID,
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

func (r *CarRepository) PublishAbortSaga(ctx context.Context,
	cmd *command.BookCar, reason string,
) error {
	db := r.WithContext(ctx)

	evt := &command.AbortSaga{
		Message: message.Message{
			Name:          reflect.ValueOf(command.AbortSaga{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			ParentID:      cmd.ParentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.AbortSagaBody{
			TripID: cmd.Body.TripID,
			Reason: reason,
			Source: "car",
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

func (r *CarRepository) PublishCarBookingCanceled(ctx context.Context,
	corrID string, parentID string, d *dto.CarBooking,
) error {
	db := r.WithContext(ctx)

	evt := &event.CarBookingCanceled{
		Message: message.Message{
			Name:          reflect.ValueOf(event.CarBookingCanceled{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarBookingCanceledBody{
			BookingID: d.ID,
			TripID:    d.TripID,
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

func (r *CarRepository) Book(ctx context.Context, d *dto.CarBooking) (dto.CarBooking, error) {
	if booking, err := r.GetByTripID(ctx, d.TripID); err != nil {
		return dto.CarBooking{}, err
	} else if booking.Status == status.Booked {
		return booking, nil
	}

	db := r.WithContext(ctx)
	row := &entity.CarBooking{
		TripID: d.TripID,
		CarID:  d.CarID,
		Status: status.Booked,
	}
	result := db.Create(row)
	if result.Error != nil {
		return dto.CarBooking{}, result.Error
	}

	return row.DTO(), nil
}

func (r *CarRepository) CancelBooking(ctx context.Context, cmd *command.CancelCarBooking) error {
	panicked := true

	tx := r.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil || panicked {
			tx.Rollback()
		}
	}()

	txCtx := context.WithValue(ctx, constant.TX("tx"), tx)

	row := &entity.CarBooking{}
	result := tx.
		Model(row).
		Where("id = ?", cmd.Body.BookingID).
		Update("status", status.Canceled)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoCarBookingFound
	}

	booking := row.DTO()
	if err := r.PublishCarBookingCanceled(txCtx, cmd.CorrelationID, cmd.ParentID, &booking); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *CarRepository) GetByID(ctx context.Context, id uint) (dto.CarBooking, error) {
	db := r.WithContext(ctx)

	row := &entity.CarBooking{}
	result := db.
		Where("id = ?", id).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.CarBooking{}, result.Error
	}
	return row.DTO(), nil
}

func (r *CarRepository) GetByTripID(ctx context.Context, tripID uint) (dto.CarBooking, error) {
	db := r.WithContext(ctx)

	row := &entity.CarBooking{}
	result := db.
		Where("trip_id = ?", tripID).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.CarBooking{}, result.Error
	}
	return row.DTO(), nil
}
