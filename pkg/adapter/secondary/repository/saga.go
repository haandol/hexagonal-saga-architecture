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
	ErrNoRowAffected = errors.New("no row affected")
)

type SagaRepository struct {
	db *gorm.DB
}

func NewSagaRepository(db *gorm.DB) *SagaRepository {
	return &SagaRepository{
		db: db,
	}
}

func (r *SagaRepository) PublishBookCar(ctx context.Context, corrID, parentID string, d dto.Saga) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	evt := &command.BookCar{
		Message: message.Message{
			Name:          reflect.ValueOf(command.BookCar{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.BookCarBody{
			TripID: d.TripID,
			CarID:  d.CarID,
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
		KafkaTopic: "car-service",
		KafkaKey:   evt.CorrelationID,
		KafkaValue: v,
	}
	return db.Create(row).Error
}

func (r *SagaRepository) PublishBookHotel(ctx context.Context, corrID, parentID string, d dto.Saga) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	evt := &command.BookHotel{
		Message: message.Message{
			Name:          reflect.ValueOf(command.BookHotel{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.BookHotelBody{
			TripID:  d.TripID,
			HotelID: d.HotelID,
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
		KafkaTopic: "hotel-service",
		KafkaKey:   evt.CorrelationID,
		KafkaValue: v,
	}
	return db.Create(row).Error
}

func (r *SagaRepository) PublishBookFlight(ctx context.Context, corrID, parentID string, d dto.Saga) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	evt := &command.BookFlight{
		Message: message.Message{
			Name:          reflect.ValueOf(command.BookFlight{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.BookFlightBody{
			TripID:   d.TripID,
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
		KafkaTopic: "flight-service",
		KafkaKey:   evt.CorrelationID,
		KafkaValue: v,
	}
	return db.Create(row).Error
}

func (r *SagaRepository) PublishEndSaga(ctx context.Context, corrID, parentID string, d dto.Saga) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	evt := &command.EndSaga{
		Message: message.Message{
			Name:          reflect.ValueOf(command.EndSaga{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.EndSagaBody{
			SagaID: d.ID,
			TripID: d.TripID,
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

func (r *SagaRepository) PublishSagaEnded(ctx context.Context, corrID, parentID string, d dto.Saga) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	evt := &event.SagaEnded{
		Message: message.Message{
			Name:          reflect.ValueOf(event.SagaEnded{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.SagaEndedBody{
			SagaID:          d.ID,
			TripID:          d.TripID,
			CarBookingID:    d.CarBookingID,
			HotelBookingID:  d.HotelBookingID,
			FlightBookingID: d.FlightBookingID,
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
		KafkaTopic: "trip-service",
		KafkaKey:   evt.CorrelationID,
		KafkaValue: v,
	}
	return db.Create(row).Error
}

func (r *SagaRepository) PublishSagaAborted(ctx context.Context, corrID, parentID string, d dto.Saga) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	evt := &event.SagaAborted{
		Message: message.Message{
			Name:          reflect.ValueOf(event.SagaAborted{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.SagaAbortedBody{
			SagaID: d.ID,
			TripID: d.TripID,
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
		KafkaTopic: "trip-service",
		KafkaKey:   evt.CorrelationID,
		KafkaValue: v,
	}
	return db.Create(row).Error
}

// get or create saga
func (r *SagaRepository) Start(ctx context.Context, cmd *command.StartSaga) error {
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

	saga, err := r.GetByTripID(txCtx, cmd.Body.TripID)
	if err != nil {
		return err
	}
	if saga.ID > 0 {
		return nil
	}
	if saga.Status == status.SagaAborted {
		return errors.New("the saga is aborting")
	}

	history, err := json.Marshal(&[]any{
		cmd,
	})
	if err != nil {
		return err
	}

	row := &entity.Saga{
		CorrelationID: cmd.CorrelationID,
		TripID:        cmd.Body.TripID,
		CarID:         cmd.Body.CarID,
		HotelID:       cmd.Body.HotelID,
		FlightID:      cmd.Body.FlightID,
		History:       history,
		Status:        status.SagaStarted,
	}
	if result := tx.Create(row); result.Error != nil {
		return result.Error
	}

	if err := r.PublishBookCar(txCtx, cmd.CorrelationID, cmd.ParentID, row.DTO()); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *SagaRepository) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error {
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

	saga, err := r.GetByCorrelationID(txCtx, evt.CorrelationID)
	if err != nil {
		return err
	}

	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return err
	}
	v = append(v, evt)

	history, err := json.Marshal(v)
	if err != nil {
		return err
	}

	row := &entity.Saga{}
	result := tx.
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"car_booking_id": evt.Body.BookingID,
			"status":         evt.Name,
			"history":        history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowAffected
	}

	if err := r.PublishBookHotel(txCtx, evt.CorrelationID, evt.ParentID, saga); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *SagaRepository) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCancelled) error {
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

	saga, err := r.GetByCorrelationID(txCtx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return err
	}
	v = append(v, evt)

	history, err := json.Marshal(v)
	if err != nil {
		return err
	}

	row := &entity.Saga{}
	result := tx.
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"car_booking_id": 0,
			"status":         status.SagaAborted,
			"history":        history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowAffected
	}

	if err := r.PublishSagaAborted(txCtx, evt.CorrelationID, evt.ParentID, saga); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *SagaRepository) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) error {
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

	saga, err := r.GetByCorrelationID(txCtx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return err
	}
	v = append(v, evt)

	history, err := json.Marshal(v)
	if err != nil {
		return err
	}

	row := &entity.Saga{}
	result := tx.
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"hotel_booking_id": evt.Body.BookingID,
			"status":           evt.Name,
			"history":          history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowAffected
	}

	if err := r.PublishBookFlight(txCtx, evt.CorrelationID, evt.ParentID, saga); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *SagaRepository) CompensateHotelBooking(ctx context.Context,
	evt *event.HotelBookingCancelled,
) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"hotel_booking_id": 0,
			"status":           status.SagaAborted,
			"history":          history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, ErrNoRowAffected
	}

	return r.GetByCorrelationID(ctx, evt.CorrelationID)
}

func (r *SagaRepository) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) error {
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

	saga, err := r.GetByCorrelationID(txCtx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return err
	}
	v = append(v, evt)

	history, err := json.Marshal(v)
	if err != nil {
		return err
	}

	row := &entity.Saga{}
	result := tx.
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"flight_booking_id": evt.Body.BookingID,
			"status":            evt.Name,
			"history":           history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowAffected
	}

	if err := r.PublishEndSaga(txCtx, evt.CorrelationID, evt.ParentID, saga); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *SagaRepository) CompensateFlightBooking(ctx context.Context,
	evt *event.FlightBookingCancelled,
) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"flight_booking_id": 0,
			"status":            status.SagaAborted,
			"history":           history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, ErrNoRowAffected
	}

	return r.GetByCorrelationID(ctx, evt.CorrelationID)
}

func (r *SagaRepository) End(ctx context.Context, cmd *command.EndSaga) error {
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

	saga, err := r.GetByCorrelationID(txCtx, cmd.CorrelationID)
	if err != nil {
		return err
	}
	if saga.Status == status.SagaEnded || saga.Status == status.SagaAborted {
		return nil
	}

	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return err
	}
	v = append(v, cmd)
	history, err := json.Marshal(v)
	if err != nil {
		return err
	}

	row := &entity.Saga{}
	result := tx.
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"status":  status.SagaEnded,
			"history": history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowAffected
	}

	if err := r.PublishSagaEnded(txCtx, cmd.CorrelationID, cmd.ParentID, saga); err != nil {
		return err
	}

	if err := tx.Commit().Error; err == nil {
		panicked = false
	} else {
		return err
	}

	return nil
}

func (r *SagaRepository) Abort(ctx context.Context, cmd *command.AbortSaga) (dto.Saga, error) {
	saga, err := r.GetByTripID(ctx, cmd.Body.TripID)
	if err != nil {
		return dto.Saga{}, err
	}
	if saga.Status == status.SagaEnded || saga.Status == status.SagaAborted {
		return saga, nil
	}

	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, cmd)
	history, err := json.Marshal(v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Table("sagas").
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"status":  status.SagaAborted,
			"history": history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, ErrNoRowAffected
	}

	return r.GetByTripID(ctx, cmd.Body.TripID)
}

func (r *SagaRepository) UpdateStatusByTripID(ctx context.Context, tripID uint, s string) error {
	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	result := db.
		Where("trip_id = ?", tripID).
		Updates(&entity.Saga{
			Status: s,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowAffected
	}

	return nil
}

func (r *SagaRepository) GetByTripID(ctx context.Context, id uint) (dto.Saga, error) {
	row := &entity.Saga{}

	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	result := db.
		Where("trip_id = ?", id).
		Limit(1).
		Find(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO(), nil
}

func (r *SagaRepository) GetByCorrelationID(ctx context.Context, id string) (dto.Saga, error) {
	row := &entity.Saga{}

	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	result := db.
		Where("correlation_id = ?", id).
		Limit(1).
		Find(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO(), nil
}
