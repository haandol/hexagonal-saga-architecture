package router

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/internal/constant"
	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/port/primaryport/routerport"
	"github.com/haandol/hexagonal/internal/service"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/cerrors"
	"github.com/haandol/hexagonal/pkg/util/o11y"
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
// @Param trip body dto.Trip true "trip id is required"
// @Success 200 {object} dto.Trip
// @Router /trips [post]
func (r *TripRouter) CreateHandler(c *gin.Context) *cerrors.CodedError {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	span := o11y.SpanFromContext(ctx)

	req := &dto.Trip{}
	if err := c.ShouldBindJSON(req); err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return cerrors.New(constant.ErrBadRequest, err)
	}

	if err := util.ValidateStruct(req); err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return cerrors.New(constant.ErrInvalidRequest, err)
	}

	trip, err := r.tripService.Create(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
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
// @Param trip_id path int true "trip id"
// @Success 200 {object} dto.Trip
// @Router /trips/{trip_id}/recover/forward [put]
func (r *TripRouter) RecoverForwardHandler(c *gin.Context) *cerrors.CodedError {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	span := o11y.SpanFromContext(ctx)

	tripID, err := strconv.Atoi(c.Param("tripId"))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return cerrors.New(constant.ErrInvalidRequest, errors.New("tripID is invalid"))
	}

	trip, err := r.tripService.RecoverForward(ctx, uint(tripID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
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
// @Param trip_id path int true "trip id"
// @Success 200 {object} dto.Trip
// @Router /trips/{trip_id}/recover/backward [put]
func (r *TripRouter) RecoverBackwardHandler(c *gin.Context) *cerrors.CodedError {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	span := o11y.SpanFromContext(ctx)

	tripID, err := strconv.Atoi(c.Param("tripId"))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return cerrors.New(constant.ErrInvalidRequest, errors.New("tripID is invalid"))
	}

	trip, err := r.tripService.RecoverBackward(ctx, uint(tripID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
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
