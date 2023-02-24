//go:build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"
	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/adapter/primary/poller"
	"github.com/haandol/hexagonal/pkg/adapter/primary/router"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/connector/database"
	kafkaconnector "github.com/haandol/hexagonal/pkg/connector/producer"
	"github.com/haandol/hexagonal/pkg/port"
	"github.com/haandol/hexagonal/pkg/port/primaryport/routerport"
	"github.com/haandol/hexagonal/pkg/service"
	"gorm.io/gorm"
)

// Common
func provideDB(cfg *config.Config) *gorm.DB {
	db, err := database.Connect(&cfg.TripDB)
	if err != nil {
		panic(err)
	}
	return db
}

func provideProducer(cfg *config.Config) *kafkaconnector.KafkaProducer {
	kafkaProducer, err := kafkaconnector.Connect(&cfg.Kafka)
	if err != nil {
		panic(err)
	}
	return kafkaProducer
}

// TripApp
func provideTripConsumer(
	cfg *config.Config,
	tripService *service.TripService,
) *consumer.TripConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "trip", "trip-service")
	return consumer.NewTripConsumer(kafkaConsumer, tripService)
}

var provideTripRouters = wire.NewSet(
	router.NewGinRouter,
	wire.Bind(new(http.Handler), new(*router.GinRouter)),
	router.NewServerForce,
	wire.Bind(new(routerport.RouterGroup), new(*router.GinRouter)),
	router.NewTripRouter,
)

var provideRouters = wire.NewSet(
	router.NewGinRouter,
	wire.Bind(new(http.Handler), new(*router.GinRouter)),
	router.NewServer,
	wire.Bind(new(routerport.RouterGroup), new(*router.GinRouter)),
)

func InitTripApp(cfg *config.Config) port.App {
	wire.Build(
		provideDB,
		provideTripRouters,
		provideTripConsumer,
		service.NewTripService,
		repository.NewTripRepository,
		NewTripApp,
		wire.Bind(new(port.App), new(*TripApp)),
	)
	return nil
}

// SagaApp
func provideSagaConsumer(
	cfg *config.Config,
	sagaService *service.SagaService,
) *consumer.SagaConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "saga", "saga-service")
	return consumer.NewSagaConsumer(kafkaConsumer, sagaService)
}

func provideSagaProducer(cfg *config.Config) *producer.SagaProducer {
	kafkaProducer := provideProducer(cfg)
	return producer.NewSagaProducer(kafkaProducer)
}

func InitSagaApp(cfg *config.Config) port.App {
	wire.Build(
		provideDB,
		provideRouters,
		provideSagaConsumer,
		service.NewSagaService,
		provideSagaProducer,
		repository.NewSagaRepository,
		NewSagaApp,
		wire.Bind(new(port.App), new(*SagaApp)),
	)
	return nil
}

// CarApp
func provideCarConsumer(
	cfg *config.Config,
	carService *service.CarService,
) *consumer.CarConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "car", "car-service")
	return consumer.NewCarConsumer(kafkaConsumer, carService)
}

func InitCarApp(cfg *config.Config) port.App {
	wire.Build(
		provideDB,
		provideRouters,
		provideCarConsumer,
		repository.NewCarRepository,
		service.NewCarService,
		NewCarApp,
		wire.Bind(new(port.App), new(*CarApp)),
	)
	return nil
}

func InitMessageRelayApp(cfg *config.Config) port.App {
	wire.Build(
		provideDB,
		provideRouters,
		poller.NewOutboxPoller,
		provideProducer,
		repository.NewOutboxRepository,
		service.NewMessageRelayService,
		NewMessageRelayApp,
		wire.Bind(new(port.App), new(*MessageRelayApp)),
	)
	return nil
}

// HotelApp
func provideHotelConsumer(
	cfg *config.Config,
	hotelService *service.HotelService,
) *consumer.HotelConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "hotel", "hotel-service")
	return consumer.NewHotelConsumer(kafkaConsumer, hotelService)
}

func InitHotelApp(cfg *config.Config) port.App {
	wire.Build(
		provideDB,
		provideRouters,
		provideHotelConsumer,
		repository.NewHotelRepository,
		service.NewHotelService,
		NewHotelApp,
		wire.Bind(new(port.App), new(*HotelApp)),
	)
	return nil
}

// FlightApp
func provideFlightConsumer(
	cfg *config.Config,
	flightService *service.FlightService,
) *consumer.FlightConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "flight", "flight-service")
	return consumer.NewFlightConsumer(kafkaConsumer, flightService)
}

func InitFlightApp(cfg *config.Config) port.App {
	wire.Build(
		provideDB,
		provideRouters,
		provideFlightConsumer,
		repository.NewFlightRepository,
		service.NewFlightService,
		NewFlightApp,
		wire.Bind(new(port.App), new(*FlightApp)),
	)
	return nil
}
