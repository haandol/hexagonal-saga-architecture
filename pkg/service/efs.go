package service

import (
	"context"
	"path/filepath"

	"github.com/haandol/hexagonal/pkg/util"
)

type EfsService struct {
}

func NewEfsService() *EfsService {
	return &EfsService{}
}

func (s *EfsService) List(ctx context.Context, path string) ([]string, error) {
	logger := util.GetLogger().WithContext(ctx).With(
		"module", "EfsService",
		"func", "List",
	)
	logger.Infow("list flies", "path", path)

	files, err := filepath.Glob(path + "/*")
	if err != nil {
		return []string{}, err
	}

	return files, nil
}
