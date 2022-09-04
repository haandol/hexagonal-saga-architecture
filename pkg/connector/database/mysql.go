package database

import (
	"context"
	"fmt"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDBs = make(map[string]*gorm.DB)

const (
	DbConnMaxLifeTime = 15 * time.Second
)

func getDsn(cfg config.Database) string {
	const postfix = "charset=utf8mb4,utf8&sql_mode=TRADITIONAL&multiStatements=true&parseTime=true&loc=Asia%2FJakarta"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, postfix)
}

func initDb(cfg config.Database) {
	if _, exists := gormDBs[cfg.Name]; exists {
		return
	}

	db, err := gorm.Open(mysql.Open(getDsn(cfg)), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	gormDBs[cfg.Name] = db
}

func Connect(cfg config.Database) (*gorm.DB, error) {
	logger := util.GetLogger()

	initDb(cfg)

	gormDB := gormDBs[cfg.Name]

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("failed to get DB instance", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(DbConnMaxLifeTime)

	logger.Infow("connected to database", "host", cfg.Host, "port", cfg.Port, "name", cfg.Name)

	return gormDB, nil
}

func Close(ctx context.Context) error {
	logger := util.GetLogger()

	var err error
	c, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		for name, db := range gormDBs {
			sqlDB, err := db.DB()
			if err != nil {
				logger.Errorw("failed to get DB instance", "name", name, "error", err)
				continue
			}

			if err := sqlDB.Close(); err != nil {
				logger.Errorw("failed to close sqlDB", "name", name, "error", err)
			}

			logger.Infow("closed database", "name", name)
		}
	}()

	<-c.Done()

	return err
}
