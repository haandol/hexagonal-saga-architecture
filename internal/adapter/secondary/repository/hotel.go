package repository

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/internal/constant"
	"github.com/haandol/hexagonal/internal/constant/status"
	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/entity"
	"github.com/haandol/hexagonal/internal/message"
	"github.com/haandol/hexagonal/internal/message/command"
	"github.com/haandol/hexagonal/internal/message/event"
	"github.com/haandol/hexagonal/pkg/util"
	"gorm.io/gorm"
)

var ErrNoHotelBookingFound = errors.New("no hotel-booking found")

type HotelRepository struct {
	BaseRepository
}

func NewHotelRepository(db *gorm.DB) *HotelRepository {
	return &HotelRepository{
		BaseRepository{DB: db},
	}
}

func (r *HotelRepository) PublishHotelBooked(ctx context.Context,
	corrID string, parentID string, d *dto.HotelBooking,
) error {
	db := r.WithContext(ctx)

	evt := &event.HotelBooked{
		Message: message.Message{
			Name:          reflect.ValueOf(event.HotelBooked{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.HotelBookedBody{
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

func (r *HotelRepository) PublishAbortSaga(ctx context.Context,
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
			Source: "hotel",
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

func (r *HotelRepository) PublishHotelBookingCanceled(ctx context.Context,
	corrID string, parentID string, d *dto.HotelBooking,
) error {
	db := r.WithContext(ctx)

	evt := &event.HotelBookingCanceled{
		Message: message.Message{
			Name:          reflect.ValueOf(event.HotelBookingCanceled{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.HotelBookingCanceledBody{
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

func (r *HotelRepository) Book(ctx context.Context, d *dto.HotelBooking, cmd *command.BookHotel) error {
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

	if booking, err := r.GetByTripID(txCtx, d.TripID); err != nil {
		return err
	} else if booking.Status == status.Booked {
		return nil
	}

	row := &entity.HotelBooking{
		TripID:  d.TripID,
		HotelID: d.HotelID,
		Status:  status.Booked,
	}
	result := tx.Create(row)
	if result.Error != nil {
		return result.Error
	}

	booking := row.DTO()
	if err := r.PublishHotelBooked(txCtx, cmd.CorrelationID, cmd.ParentID, &booking); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *HotelRepository) CancelBooking(ctx context.Context, cmd *command.CancelHotelBooking) error {
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

	row := &entity.HotelBooking{}
	result := tx.
		Model(row).
		Where("id = ?", cmd.Body.BookingID).
		Update("status", status.Canceled)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoHotelBookingFound
	}

	booking := row.DTO()
	if err := r.PublishHotelBookingCanceled(txCtx, cmd.CorrelationID, cmd.ParentID, &booking); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *HotelRepository) GetByID(ctx context.Context, id uint) (dto.HotelBooking, error) {
	db := r.WithContext(ctx)

	row := &entity.HotelBooking{}
	result := db.
		Where("id = ?", id).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}
	return row.DTO(), nil
}

func (r *HotelRepository) GetByTripID(ctx context.Context, tripID uint) (dto.HotelBooking, error) {
	db := r.WithContext(ctx)

	row := &entity.HotelBooking{}
	result := db.
		Where("trip_id = ?", tripID).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}
	return row.DTO(), nil
}
