package config

import (
	"log"

	"github.com/haandol/hexagonal/pkg/util"
	"github.com/joho/godotenv"
)

type App struct {
	Name                    string `validate:"required"`
	Stage                   string `validate:"required"`
	Port                    int    `validate:"required,number"`
	RPS                     int    `validate:"required,number"`
	TimeoutSec              int    `validate:"required,number,gte=0"`
	GracefulShutdownTimeout int    `validate:"required,number,gte=0"`
	DisableHTTP             bool   `default:"false"`
}

type Kafka struct {
	Seeds            []string `validate:"required"`
	MessageExpirySec int      `validate:"required,number"`
	BatchSize        int      `validate:"required,number"`
}

type Database struct {
	Host               string
	Port               int
	Name               string
	Username           string
	Password           string
	SecretID           string
	MaxOpenConnections int `validate:"required,number"`
	MaxIdleConnections int `validate:"required,number"`
}

type Relay struct {
	FetchSize        int `validate:"required,number"`
	FetchIntervalMil int `validate:"required,number"`
}

type Config struct {
	App    App
	Kafka  Kafka
	TripDB Database
	Relay  Relay
}

// Load config.Config from environment variables for each stage
func Load() Config {
	stage := getEnv("APP_STAGE").String()
	log.Printf("Loading %s config\n", stage)

	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	cfg := Config{
		App: App{
			Name:                    getEnv("APP_NAME").String(),
			Stage:                   getEnv("APP_STAGE").String(),
			Port:                    getEnv("APP_PORT").Int(),
			RPS:                     getEnv("APP_RPS").Int(),
			TimeoutSec:              getEnv("APP_TIMEOUT_SEC").Int(),
			GracefulShutdownTimeout: getEnv("APP_GRACEFUL_SHUTDOWN_TIMEOUT").Int(),
			DisableHTTP:             getEnv("APP_DISABLE_HTTP").Bool(),
		},
		Kafka: Kafka{
			Seeds:            getEnv("KAFKA_SEEDS").Split(","),
			MessageExpirySec: getEnv("KAFKA_MESSAGE_EXPIRY_SEC").Int(),
			BatchSize:        getEnv("KAFKA_BATCH_SIZE").Int(),
		},
		TripDB: Database{
			Host:               getEnv("DB_HOST").String(),
			Port:               getEnv("DB_PORT").Int(),
			Name:               getEnv("DB_NAME").String(),
			Username:           getEnv("DB_USERNAME").String(),
			Password:           getEnv("DB_PASSWORD").String(),
			SecretID:           getEnv("DB_SECRET_ID").String(),
			MaxOpenConnections: getEnv("DB_MAX_OPEN_CONNECTIONS").Int(),
			MaxIdleConnections: getEnv("DB_MAX_IDLE_CONNECTIONS").Int(),
		},
		Relay: Relay{
			FetchSize:        getEnv("RELAY_FETCH_SIZE").Int(),
			FetchIntervalMil: getEnv("RELAY_FETCH_INTERVAL_MIL").Int(),
		},
	}

	if err := util.ValidateStruct(cfg); err != nil {
		log.Panicf("Error validating config: %s", err.Error())
	}

	if cfg.TripDB.SecretID == "" {
		if err := util.ValidateVar(cfg.TripDB.Host, "required"); err != nil {
			log.Panicf("Error validating config: %s", err.Error())
		}
		if err := util.ValidateVar(cfg.TripDB.Port, "required"); err != nil {
			log.Panicf("Error validating config: %s", err.Error())
		}
		if err := util.ValidateVar(cfg.TripDB.Username, "required"); err != nil {
			log.Panicf("Error validating config: %s", err.Error())
		}
		if err := util.ValidateVar(cfg.TripDB.Password, "required"); err != nil {
			log.Panicf("Error validating config: %s", err.Error())
		}
	}

	return cfg
}
