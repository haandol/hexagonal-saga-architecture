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

func (r *OutboxRepository) QueryUnsent(ctx context.Context, batchSize int) ([]dto.Outbox, error) {
	rows := entity.Outboxes{}
	result := r.db.WithContext(ctx).
		Where("is_sent = ?", false).
		Limit(batchSize).
		Order("id ASC").
		Find(&rows)
	if result.Error != nil {
		return []dto.Outbox{}, result.Error
	}

	return rows.DTO(), nil
}

func (r *OutboxRepository) MarkSent(ctx context.Context, id uint) error {
	row := entity.Outbox{
		IsSent: true,
	}
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(&row).Error
}
