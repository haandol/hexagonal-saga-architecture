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
	ErrNoFlightBookingFound = errors.New("no flight-booking found")
)

type FlightRepository struct {
	BaseRepository
}

func NewFlightRepository(db *gorm.DB) *FlightRepository {
	return &FlightRepository{
		BaseRepository{DB: db},
	}
}

func (r *FlightRepository) PublishFlightBooked(ctx context.Context,
	corrID string, parentID string, d *dto.FlightBooking,
) error {
	db := r.WithContext(ctx)

	evt := &event.FlightBooked{
		Message: message.Message{
			Name:          reflect.ValueOf(event.FlightBooked{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.FlightBookedBody{
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

func (r *FlightRepository) PublishAbortSaga(ctx context.Context,
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
			Source: "flight",
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

func (r *FlightRepository) PublishFlightBookingCancelled(ctx context.Context,
	corrID string, parentID string, d *dto.FlightBooking,
) error {
	db := r.WithContext(ctx)

	evt := &event.FlightBookingCancelled{
		Message: message.Message{
			Name:          reflect.ValueOf(event.FlightBookingCancelled{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.FlightBookingCancelledBody{
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

func (r *FlightRepository) Book(ctx context.Context, d *dto.FlightBooking, cmd *command.BookFlight) error {
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

	row := &entity.FlightBooking{
		TripID:   d.TripID,
		FlightID: d.FlightID,
		Status:   status.Booked,
	}
	result := tx.Create(row)
	if result.Error != nil {
		return result.Error
	}

	booking := row.DTO()
	if err := r.PublishFlightBooked(txCtx, cmd.CorrelationID, cmd.ParentID, &booking); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *FlightRepository) CancelBooking(ctx context.Context, cmd *command.CancelFlightBooking) error {
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

	row := &entity.FlightBooking{}
	result := tx.
		Model(row).
		Where("id = ?", cmd.Body.BookingID).
		Update("status", status.Cancelled)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoFlightBookingFound
	}

	booking := row.DTO()
	if err := r.PublishFlightBookingCancelled(txCtx, cmd.CorrelationID, cmd.ParentID, &booking); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *FlightRepository) GetByID(ctx context.Context, id uint) (dto.FlightBooking, error) {
	row := &entity.FlightBooking{}
	result := r.WithContext(ctx).
		Where("id = ?", id).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.FlightBooking{}, result.Error
	}
	return row.DTO(), nil
}

func (r *FlightRepository) GetByTripID(ctx context.Context, tripID uint) (dto.FlightBooking, error) {
	db := r.WithContext(ctx)

	row := &entity.FlightBooking{}
	result := db.
		Where("trip_id = ?", tripID).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.FlightBooking{}, result.Error
	}
	return row.DTO(), nil
}
