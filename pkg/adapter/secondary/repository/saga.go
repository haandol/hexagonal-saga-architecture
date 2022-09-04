package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"gorm.io/gorm"
)

type SagaRepository struct {
	db *gorm.DB
}

func NewSagaRepository(db *gorm.DB) *SagaRepository {
	return &SagaRepository{
		db: db,
	}
}

func (r *SagaRepository) Start(ctx context.Context, cmd *command.StartSaga) (dto.Saga, error) {
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
		Status:        "STARTED",
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO()
}

func (r *SagaRepository) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error {
	saga, err := r.GetByCorrelationId(ctx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, evt)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			CarBookingID: evt.Body.BookingID,
			Status:       evt.Name,
			History:      history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCanceled) error {
	saga, err := r.GetByCorrelationId(ctx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, evt)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			CarBookingID: 0,
			Status:       evt.Name,
			History:      history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) error {
	saga, err := r.GetByCorrelationId(ctx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, evt)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			HotelBookingID: evt.Body.BookingID,
			Status:         evt.Name,
			History:        history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCanceled) error {
	saga, err := r.GetByCorrelationId(ctx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, evt)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			HotelBookingID: 0,
			Status:         evt.Name,
			History:        history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) error {
	saga, err := r.GetByCorrelationId(ctx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, evt)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			FlightBookingID: evt.Body.BookingID,
			Status:          evt.Name,
			History:         history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCanceled) error {
	saga, err := r.GetByCorrelationId(ctx, evt.CorrelationID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, evt)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			FlightBookingID: 0,
			Status:          evt.Name,
			History:         history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) End(ctx context.Context, cmd *command.EndSaga) error {
	saga, err := r.GetById(ctx, cmd.Body.SagaID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, cmd)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", cmd.Body.SagaID).
		Updates(&entity.Saga{
			Status:  "ENDED",
			History: history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) Abort(ctx context.Context, cmd *command.AbortSaga) error {
	saga, err := r.GetById(ctx, cmd.Body.SagaID)
	if err != nil {
		return err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &[]any{}); err != nil {
		return err
	}
	v = append(v, cmd)
	history, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).
		Where("id = ?", cmd.Body.SagaID).
		Updates(&entity.Saga{
			Status:  "ABORTED",
			History: history,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *SagaRepository) GetById(ctx context.Context, id uint) (dto.Saga, error) {
	row := &entity.Saga{}

	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	result := db.
		Where("id = ?", id).
		Take(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO()
}

func (r *SagaRepository) GetByCorrelationId(ctx context.Context, id string) (dto.Saga, error) {
	row := &entity.Saga{}

	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	result := db.
		Where("correlation_id = ?", id).
		Find(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO()
}
