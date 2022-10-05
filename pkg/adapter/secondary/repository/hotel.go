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
	ErrNoHotelBookingFound = errors.New("not hotel-booking found")
)

type HotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) *HotelRepository {
	return &HotelRepository{
		db: db,
	}
}

func (r *HotelRepository) PublishHotelBooked(ctx context.Context,
	corrID string, parentID string, d dto.HotelBooking,
) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

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
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

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

func (r *HotelRepository) PublishHotelBookingCancelled(ctx context.Context,
	corrID string, parentID string, d dto.HotelBooking,
) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	evt := &event.HotelBookingCancelled{
		Message: message.Message{
			Name:          reflect.ValueOf(event.HotelBookingCancelled{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.HotelBookingCancelledBody{
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

	tx := r.db.WithContext(ctx).Begin()
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

	if err := r.PublishHotelBooked(txCtx, cmd.CorrelationID, cmd.ParentID, row.DTO()); err != nil {
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

	tx := r.db.WithContext(ctx).Begin()
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
		Update("status", status.Cancelled)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoHotelBookingFound
	}

	if err := r.PublishHotelBookingCancelled(txCtx, cmd.CorrelationID, cmd.ParentID, row.DTO()); err != nil {
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
	row := &entity.HotelBooking{}
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}
	return row.DTO(), nil
}

func (r *HotelRepository) GetByTripID(ctx context.Context, tripID uint) (dto.HotelBooking, error) {
	row := &entity.HotelBooking{}
	result := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}
	return row.DTO(), nil
}
