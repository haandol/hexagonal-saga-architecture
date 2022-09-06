package router

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	g.Handle("PUT", "/:tripId/recover/forward", r.WrappedHandler(r.RecoverForwardHandler))
	g.Handle("PUT", "/:tripId/recover/backward", r.WrappedHandler(r.RecoverBackwardHandler))
}

// @Summary create new trip
// @Schemes
// @Description create new trip
// @Tags trips
// @Accept json
// @Produce json
// @Param "trip" body dto.Trip true "trip id is required"
// @Success 200 {object} dto.Trip
// @Router /trips [post]
func (r *TripRouter) CreateHandler(c *gin.Context) *cerrors.CodedError {
	req := &dto.Trip{}
	if err := c.ShouldBindJSON(req); err != nil {
		return cerrors.New(constant.ErrBadRequest, err)
	}

	corrID := c.Request.Header.Get("X-Request-ID")
	if corrID == "" {
		corrID = uuid.NewString()
	}

	if err := util.ValidateStruct(req); err != nil {
		return cerrors.New(constant.ErrInvalidRequest, err)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	trip, err := r.tripService.Create(ctx, corrID, req)
	if err != nil {
		return cerrors.New(constant.ErrFailToCreateTrip, err)
	}

	return r.Success(c, trip)
}

// @Summary recover forward
// @Schemes
// @Description recover forward
// @Tags trips
// @Accept json
// @Produce json
// @Param "tripId" path uint true "trip id"
// @Success 200 {object} dto.Trip
// @Router /trips/{tripId}/recover/forward [put]
func (r *TripRouter) RecoverForwardHandler(c *gin.Context) *cerrors.CodedError {
	corrID := c.Request.Header.Get("X-Request-ID")
	if corrID == "" {
		return cerrors.New(constant.ErrBadRequest, errors.New("X-Request-ID is required"))
	}

	tripID, err := strconv.ParseUint(c.Param("tripId"), 10, 32)
	if err != nil {
		return cerrors.New(constant.ErrInvalidRequest, errors.New("tripID is invalid"))
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	trip, err := r.tripService.RecoverForward(ctx, corrID, uint(tripID))
	if err != nil {
		return cerrors.New(constant.ErrFailToCreateTrip, err)
	}

	return r.Success(c, trip)
}

// @Summary recover backward
// @Schemes
// @Description recover backward
// @Tags trips
// @Accept json
// @Produce json
// @Param "tripId" path uint true "trip id"
// @Success 200 {object} dto.Trip
// @Router /trips/{tripId}/recover/backward [put]
func (r *TripRouter) RecoverBackwardHandler(c *gin.Context) *cerrors.CodedError {
	corrID := c.Request.Header.Get("X-Request-ID")
	if corrID == "" {
		return cerrors.New(constant.ErrBadRequest, errors.New("X-Request-ID is required"))
	}

	tripID, err := strconv.ParseUint(c.Param("tripId"), 10, 32)
	if err != nil {
		return cerrors.New(constant.ErrInvalidRequest, errors.New("tripID is invalid"))
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	trip, err := r.tripService.RecoverBackward(ctx, corrID, uint(tripID))
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
func (r *TripRouter) ListHandler(c *gin.Context) *cerrors.CodedError {
	trips, err := r.tripService.List(c.Request.Context())
	if err != nil {
		return cerrors.New(constant.ErrFailToListTrip, err)
	}

	return r.Success(c, trips)
}
