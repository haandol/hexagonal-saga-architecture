//go:build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"
	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/adapter/primary/router"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/connector/database"
	"github.com/haandol/hexagonal/pkg/port"
	"github.com/haandol/hexagonal/pkg/port/primaryport/routerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/service"
	"gorm.io/gorm"
)

// TripApp
func provideTripDB(cfg config.Config) *gorm.DB {
	db, err := database.Connect(cfg.TripDB)
	if err != nil {
		panic(err)
	}
	return db
}

var provideProducer = wire.NewSet(
	producer.NewKafkaProducer,
	wire.Bind(new(producerport.Producer), new(*producer.KafkaProducer)),
)

var provideRepositories = wire.NewSet(
	repository.NewTripRepository,
	wire.Bind(new(repositoryport.TripRepository), new(*repository.TripRepository)),
)

var provideRestServices = wire.NewSet(
	service.NewTripService,
)

var provideRouters = wire.NewSet(
	router.NewGinRouter,
	wire.Bind(new(http.Handler), new(*router.GinRouter)),
	wire.Bind(new(routerport.RouterGroup), new(*router.GinRouter)),
	router.NewTripRouter,
)

func InitTripApp(cfg config.Config) port.App {
	wire.Build(
		provideTripDB,
		provideRepositories,
		provideRestServices,
		provideRouters,
		provideProducer,
		NewServer,
		NewTripApp,
		wire.Bind(new(port.App), new(*TripApp)),
	)
	return nil
}

// SagaApp
func provideSagaConsumer(
	cfg config.Config,
	sagaService *service.SagaService,
) *consumer.SagaConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(cfg.Kafka, "saga", "saga-service")
	return consumer.NewSagaConsumer(kafkaConsumer, sagaService)
}

func provideSagaProducer(cfg config.Config) *producer.SagaProducer {
	kafkaProducer := producer.NewKafkaProducer(cfg)
	return producer.NewSagaProducer(kafkaProducer)
}

func InitSagaApp(cfg config.Config) port.App {
	wire.Build(
		provideTripDB,
		provideSagaConsumer,
		service.NewSagaService,
		provideSagaProducer,
		wire.Bind(new(producerport.SagaProducer), new(*producer.SagaProducer)),
		repository.NewSagaRepository,
		wire.Bind(new(repositoryport.SagaRepository), new(*repository.SagaRepository)),
		NewSagaApp,
		wire.Bind(new(port.App), new(*SagaApp)),
	)
	return nil
}

// CarApp
func provideCarConsumer(
	cfg config.Config,
	carService *service.CarService,
) *consumer.CarConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(cfg.Kafka, "car", "car-service")
	return consumer.NewCarConsumer(kafkaConsumer, carService)
}

func provideCarProducer(cfg config.Config) *producer.CarProducer {
	kafkaProducer := producer.NewKafkaProducer(cfg)
	return producer.NewCarProducer(kafkaProducer)
}

func InitCarApp(cfg config.Config) port.App {
	wire.Build(
		provideTripDB,
		provideCarConsumer,
		service.NewCarService,
		provideCarProducer,
		wire.Bind(new(producerport.CarProducer), new(*producer.CarProducer)),
		repository.NewCarRepository,
		wire.Bind(new(repositoryport.CarRepository), new(*repository.CarRepository)),
		NewCarApp,
		wire.Bind(new(port.App), new(*CarApp)),
	)
	return nil
}

// HotelApp
func provideHotelConsumer(
	cfg config.Config,
	hotelService *service.HotelService,
) *consumer.HotelConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(cfg.Kafka, "hotel", "hotel-service")
	return consumer.NewHotelConsumer(kafkaConsumer, hotelService)
}

func provideHotelProducer(cfg config.Config) *producer.HotelProducer {
	kafkaProducer := producer.NewKafkaProducer(cfg)
	return producer.NewHotelProducer(kafkaProducer)
}

func InitHotelApp(cfg config.Config) port.App {
	wire.Build(
		provideTripDB,
		provideHotelConsumer,
		service.NewHotelService,
		repository.NewHotelRepository,
		wire.Bind(new(repositoryport.HotelRepository), new(*repository.HotelRepository)),
		provideHotelProducer,
		wire.Bind(new(producerport.HotelProducer), new(*producer.HotelProducer)),
		NewHotelApp,
		wire.Bind(new(port.App), new(*HotelApp)),
	)
	return nil
}

// FlightApp
func provideFlightConsumer(
	cfg config.Config,
	flightService *service.FlightService,
) *consumer.FlightConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(cfg.Kafka, "flight", "flight-service")
	return consumer.NewFlightConsumer(kafkaConsumer, flightService)
}

func provideFlightProducer(cfg config.Config) *producer.FlightProducer {
	kafkaProducer := producer.NewKafkaProducer(cfg)
	return producer.NewFlightProducer(kafkaProducer)
}

func InitFlightApp(cfg config.Config) port.App {
	wire.Build(
		provideTripDB,
		provideFlightConsumer,
		service.NewFlightService,
		repository.NewFlightRepository,
		wire.Bind(new(repositoryport.FlightRepository), new(*repository.FlightRepository)),
		provideFlightProducer,
		wire.Bind(new(producerport.FlightProducer), new(*producer.FlightProducer)),
		NewFlightApp,
		wire.Bind(new(port.App), new(*FlightApp)),
	)
	return nil
}
