package repository

import (
	"context"

	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/entity"
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

func (r *OutboxRepository) MarkSentInBatch(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).
		Model(&entity.Outbox{}).
		Where("id IN ?", ids).
		UpdateColumn("is_sent", true).Error
}
