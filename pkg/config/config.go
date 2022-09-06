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
}

type Kafka struct {
	Seeds            []string `validate:"required"`
	MessageExpirySec int      `validate:"required,number"`
	BatchSize        int      `validate:"required,number"`
}

type Database struct {
	Host               string `validate:"required"`
	Port               int    `validate:"required,number"`
	Name               string `validate:"required"`
	Username           string `validate:"required"`
	Password           string `validate:"required"`
	MaxOpenConnections int    `validate:"required,number"`
	MaxIdleConnections int    `validate:"required,number"`
}

type Config struct {
	App    App
	Kafka  Kafka
	TripDB Database
}

func BuildFromPath(envPath string) Config {
	if err := godotenv.Load(envPath); err != nil {
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
			MaxOpenConnections: getEnv("DB_MAX_OPEN_CONNECTIONS").Int(),
			MaxIdleConnections: getEnv("DB_MAX_IDLE_CONNECTIONS").Int(),
		},
	}

	if err := util.ValidateStruct(cfg); err != nil {
		log.Panicf("Error validating config: %s", err.Error())
	}
	return cfg
}

// Load config.Config from environment variables for each stage
func Load() Config {
	stage := getEnv("APP_STAGE").String()
	log.Printf("Loading %s config\n", stage)

	envPath := getEnv("DOTENV_PATH").String()
	// use local.env for the only dev env
	if stage == "" && envPath == "" {
		return BuildFromPath("../../env/local.env")
	}

	return BuildFromPath(envPath)
}
