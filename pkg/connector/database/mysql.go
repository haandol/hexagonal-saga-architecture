package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/util"
	"moul.io/zapgorm2"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDBs = make(map[string]*gorm.DB)

const (
	DBConnMaxLifeTime = 15 * time.Second
)

func getSQLDsn(cfg *config.Database) string {
	const postfix = "charset=utf8mb4,utf8&sql_mode=TRADITIONAL&multiStatements=true&parseTime=true&loc=Asia%2FSeoul"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, postfix)
}

func initDB(cfg *config.Database) {
	if _, exists := gormDBs[cfg.Name]; exists {
		return
	}

	logger := zapgorm2.New(util.GetLogger().With("package", "mysql").Desugar())
	logger.SetAsDefault()
	db, err := gorm.Open(
		mysql.Open(getSQLDsn(cfg)),
		&gorm.Config{
			Logger:      logger,
			PrepareStmt: false,
		},
	)
	if err != nil {
		panic(err)
	}

	gormDBs[cfg.Name] = db
}

func Connect(cfg *config.Database) (*gorm.DB, error) {
	logger := util.GetLogger().With(
		"pkg", "database",
		"func", "Connect",
	)

	initDB(cfg)

	gormDB := gormDBs[cfg.Name]

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("failed to get DB instance", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(DBConnMaxLifeTime)

	logger.Infow("Connected to database", "host", cfg.Host, "port", cfg.Port, "name", cfg.Name)

	return gormDB, nil
}

func Close(ctx context.Context) error {
	logger := util.GetLogger().With(
		"pkg", "database",
		"func", "Close",
	)
	logger.Info("Closing database connection...")

	var sqlDB *sql.DB
	var err error

	done := make(chan error)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		for name, db := range gormDBs {
			sqlDB, err = db.DB()
			if err != nil {
				logger.Errorw("failed to get DB instance", "name", name, "error", err)
				continue
			}

			if err = sqlDB.Close(); err != nil {
				logger.Errorw("failed to close sqlDB", "name", name, "error", err)
			}

			logger.Infow("Closed database", "name", name)
		}
		done <- nil
	}()

	select {
	case <-done:
		return err
	case <-ctx.Done():
		return errors.New("timeout closing redis connection")
	}
}
