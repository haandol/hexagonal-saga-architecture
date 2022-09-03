package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/util"
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

func (r *SagaRepository) End(ctx context.Context, cmd *command.EndSaga) error {
	logger := util.GetLogger().With(
		"pkg", "repository",
		"module", "SagaRepository",
		"func", "End",
	)

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Rollback")
			tx.Rollback()
		}
	}()

	txCtx := context.WithValue(ctx, constant.TX("tx"), tx)

	saga, err := r.GetById(txCtx, cmd.Body.SagaID)
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

	result := tx.
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
	logger := util.GetLogger().With(
		"pkg", "repository",
		"module", "SagaRepository",
		"func", "Abort",
	)

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Rollback")
			tx.Rollback()
		}
	}()

	txCtx := context.WithValue(ctx, constant.TX("tx"), tx)

	saga, err := r.GetById(txCtx, cmd.Body.SagaID)
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

	result := tx.
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

	tx.Commit()
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
