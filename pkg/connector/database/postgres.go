package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gormDBs = make(map[string]*gorm.DB)

const (
	DBConnMaxLifeTime = 15 * time.Second
)

func getSQLDsn(cfg config.Database) string {
	const dsn = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
	return fmt.Sprintf(dsn, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}

func initDB(cfg config.Database) {
	if _, exists := gormDBs[cfg.Name]; exists {
		return
	}

	fmt.Println("dsn: ", getSQLDsn(cfg))
	instrumentedDB, err := xray.SQLContext("postgres", getSQLDsn(cfg))
	if err != nil {
		fmt.Println("!!!!!!!!!!")
		panic(err)
	}

	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: instrumentedDB}),
		&gorm.Config{PrepareStmt: true},
	)
	if err != nil {
		panic(err)
	}

	gormDBs[cfg.Name] = db
}

func Connect(cfg config.Database) (*gorm.DB, error) {
	logger := util.GetLogger()

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
