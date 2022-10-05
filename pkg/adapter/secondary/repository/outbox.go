package repository

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"gorm.io/gorm"
)

type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{
		db: db,
	}
}

func (r *OutboxRepository) QueryUnsent(ctx context.Context) ([]dto.Outbox, error) {
	rows := entity.Outboxes{}
	result := r.db.WithContext(ctx).
		Limit(100).
		Order("id ASC").
		Find(&rows)
	if result.Error != nil {
		return []dto.Outbox{}, result.Error
	}

	return rows.DTO(), nil
}

func (r *OutboxRepository) Delete(ctx context.Context, id uint) error {
	row := entity.Outbox{
		ID: id,
	}
	return r.db.WithContext(ctx).
		Unscoped().
		Delete(&row).Error
}
