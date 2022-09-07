package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/constant/status"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SagaRepository struct {
	db *gorm.DB
}

func NewSagaRepository(db *gorm.DB) *SagaRepository {
	return &SagaRepository{
		db: db,
	}
}

// get or create saga
func (r *SagaRepository) Start(ctx context.Context, cmd *command.StartSaga) (dto.Saga, error) {
	saga, err := r.GetByTripID(ctx, cmd.Body.TripID)
	if err != nil {
		return dto.Saga{}, err
	}
	if saga.ID > 0 {
		return saga, nil
	}

	history, err := json.Marshal(&[]any{
		cmd,
	})
	if err != nil {
		return dto.Saga{}, err
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
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO()
}

func (r *SagaRepository) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) (dto.Saga, error) {
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
		Clauses(clause.Returning{}).
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"car_booking_id": evt.Body.BookingID,
			"status":         evt.Name,
			"history":        history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCancelled) (dto.Saga, error) {
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
		Clauses(clause.Returning{}).
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"car_booking_id": 0,
			"status":         evt.Name,
			"history":        history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) (dto.Saga, error) {
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
		Clauses(clause.Returning{}).
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"hotel_booking_id": evt.Body.BookingID,
			"status":           evt.Name,
			"history":          history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCancelled) (dto.Saga, error) {
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
		Clauses(clause.Returning{}).
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"hotel_booking_id": 0,
			"status":           evt.Name,
			"history":          history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) (dto.Saga, error) {
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
		Clauses(clause.Returning{}).
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"flight_booking_id": evt.Body.BookingID,
			"status":            evt.Name,
			"history":           history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCancelled) (dto.Saga, error) {
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
		Clauses(clause.Returning{}).
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"flight_booking_id": 0,
			"status":            evt.Name,
			"history":           history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) End(ctx context.Context, cmd *command.EndSaga) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, cmd.CorrelationID)
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
		Clauses(clause.Returning{}).
		Limit(1).
		Where("id = ?", saga.ID).
		Updates(map[string]interface{}{
			"status":  status.SagaEnded,
			"history": history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
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
		Clauses(clause.Returning{}).
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
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) UpdateStatusByTripID(ctx context.Context, tripID uint, s string) error {
	result := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Updates(&entity.Saga{
			Status: s,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
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

	return row.DTO()
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

	return row.DTO()
}
