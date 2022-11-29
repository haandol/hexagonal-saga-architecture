package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/connector/cloud"
	"github.com/haandol/hexagonal/pkg/connector/database/internal"
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
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Name, postfix,
	)
}

func getSQLDsnFromSecretsManager(cfg *config.Database) string {
	awsCfg, err := cloud.GetAWSConfigWithProfile("skt")
	if err != nil {
		log.Fatalf("unable to get AWS Config %s", err.Error())
	}

	secrets, err := internal.GetSecretsWithID(&awsCfg.Cfg, cfg.SecretID)
	if err != nil {
		log.Fatalf("unable to get secrets %s", err.Error())
	}

	const postfix = "charset=utf8mb4,utf8&sql_mode=TRADITIONAL&multiStatements=true&parseTime=true&loc=Asia%2FSeoul"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		secrets.Username, secrets.Password,
		secrets.Host, secrets.Port,
		cfg.Name, postfix,
	)
}

func initDB(cfg *config.Database) {
	if _, exists := gormDBs[cfg.Name]; exists {
		return
	}

	logger := zapgorm2.New(util.GetLogger().With("package", "mysql").Desugar())
	logger.SetAsDefault()

	var dsn string
	if cfg.SecretID == "" {
		log.Println("dsn from config")
		dsn = getSQLDsn(cfg)
	} else {
		log.Println("dsn from secrets")
		dsn = getSQLDsnFromSecretsManager(cfg)
	}

	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger:      logger,
			PrepareStmt: false,
		},
	)
	if err != nil {
		log.Fatalf("failed to open gorm %s", err.Error())
	}

	/* TODO: remove this if you want to trace db query
	if err := db.Use(internal.NewPlugin()); err != nil {
		log.Fatalf("failed to use otel plugin for gorm: %s", err.Error())
	}
	*/

	gormDBs[cfg.Name] = db
}

func Connect(cfg *config.Database) (*gorm.DB, error) {
	logger := util.GetLogger().With(
		"pkg", "database",
		"func", "Connect",
	)
	logger.Infow("Connecting to database...")

	initDB(cfg)

	gormDB := gormDBs[cfg.Name]

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("failed to get DB instance", err.Error())
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(DBConnMaxLifeTime)

	logger.Infow("Connected to database")

	return gormDB, nil
}

func Close(ctx context.Context) error {
	logger := util.GetLogger().With(
		"pkg", "database",
		"func", "Close",
	)
	logger.Info("Closing database connection...")

	done := make(chan error)
	go func() {
		for name, db := range gormDBs {
			sqlDB, err := db.DB()
			if err != nil {
				logger.Errorw("failed to get DB instance", "name", name, "error", err.Error())
				continue
			}

			if err = sqlDB.Close(); err != nil {
				logger.Errorw("failed to close sqlDB", "name", name, "error", err.Error())
			}

			logger.Infow("Closed database", "name", name)
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return errors.New("timeout closing redis connection")
	}
}
