//go:build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"
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
		provideProducer,
		provideRepositories,
		provideRestServices,
		provideRouters,
		NewServer,
		NewTripApp,
		wire.Bind(new(port.App), new(*TripApp)),
	)
	return nil
}
