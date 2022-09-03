package router

import (
	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/port/primaryport/routerport"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/cerrors"
)

type TripRouter struct {
	BaseRouter
	tripService *service.TripService
}

func NewTripRouter(
	tripService *service.TripService,
) *TripRouter {
	return &TripRouter{
		tripService: tripService,
	}
}

func (r *TripRouter) Route(rg routerport.RouterGroup) {
	g := rg.Group("/trips")
	g.Handle("POST", "/", r.WrappedHandler(r.CreateHandler))
	g.Handle("GET", "/", r.WrappedHandler(r.ListHandler))
}

// @Summary create new trip
// @Schemes
// @Description create new trip
// @Tags trips
// @Accept json
// @Produce json
// @Param "userId" body uint true "user id"
// @Success 200 {object} dto.Trip
// @Router /trips [post]
// @Security BearerAuth
func (r *TripRouter) CreateHandler(c *gin.Context) *cerrors.CodedError {
	req := &dto.Trip{}
	if err := c.ShouldBindJSON(req); err != nil {
		return cerrors.New(constant.ErrBadRequest, err)
	}

	if err := util.ValidateStruct(req); err != nil {
		return cerrors.New(constant.ErrInvalidRequest, err)
	}

	trip, err := r.tripService.Create(c.Request.Context(), req)
	if err != nil {
		return cerrors.New(constant.ErrFailToCreateTrip, err)
	}

	return r.Success(c, trip)
}

// @Summary list all trips
// @Schemes
// @Description list all trips
// @Tags trips
// @Accept json
// @Produce json
// @Success 200 {object} []dto.Trip
// @Router /trips [get]
// @Security BearerAuth
func (r *TripRouter) ListHandler(c *gin.Context) *cerrors.CodedError {
	trips, err := r.tripService.List(c.Request.Context())
	if err != nil {
		return cerrors.New(constant.ErrFailToListTrip, err)
	}

	return r.Success(c, trips)
}
