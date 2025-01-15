package repository

import (
	"context"

	"github.com/pkg/errors"

	"github.com/haandol/hexagonal/internal/constant"
	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func (r BaseRepository) WithContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		return tx
	} else {
		return r.DB.WithContext(ctx)
	}
}

func (r BaseRepository) BeginTx(ctx context.Context) (context.Context, error) {
	if _, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		return ctx, errors.New("transaction already exists")
	}

	tx := r.DB.Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}

	return context.WithValue(ctx, constant.TX("tx"), tx), nil
}

func (r BaseRepository) CommitTx(ctx context.Context) error {
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		return tx.Commit().Error
	}

	return errors.New("no transaction found")
}

func (r BaseRepository) RollbackTx(ctx context.Context) error {
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		return tx.Rollback().Error
	}

	return errors.New("no transaction found")
}
