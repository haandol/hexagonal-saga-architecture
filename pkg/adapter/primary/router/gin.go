package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/docs"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/middleware"
	"github.com/haandol/hexagonal/pkg/port/primaryport/routerport"
	"github.com/haandol/hexagonal/pkg/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinRouter struct {
	*gin.Engine
}

func (r *GinRouter) Use(middlewares ...any) {
	for _, mw := range middlewares {
		h := mw.(func(*gin.Context))
		r.Engine.Use(h)
	}
}

func (r *GinRouter) Group(path string) routerport.RouterGroup {
	return &GinRouterGroup{
		r.Engine.Group(path),
	}
}

func (r *GinRouter) Handle(method, path string, handlerFunc ...any) {
	var ginHandlers []gin.HandlerFunc
	for _, handler := range handlerFunc {
		h := handler.(func(*gin.Context))
		ginHandlers = append(ginHandlers, h)
	}

	r.Engine.Handle(method, path, ginHandlers...)
}

type GinRouterGroup struct {
	*gin.RouterGroup
}

func (r *GinRouterGroup) Use(middlewares ...any) {
	for _, mw := range middlewares {
		h := getHandlerFunc(mw)
		r.RouterGroup.Use(h)
	}
}

func (r *GinRouterGroup) Group(path string) routerport.RouterGroup {
	return &GinRouterGroup{
		r.RouterGroup.Group(path),
	}
}

func (r *GinRouterGroup) Handle(method, path string, handlerFunc ...any) {
	var ginHandlers []gin.HandlerFunc
	for _, handler := range handlerFunc {
		h := getHandlerFunc(handler)
		ginHandlers = append(ginHandlers, h)
	}

	r.RouterGroup.Handle(method, path, ginHandlers...)
}

// fun(*gin.Context) or gin.HandlerFunc
func getHandlerFunc(f any) gin.HandlerFunc {
	if h, ok := f.(func(*gin.Context)); !ok {
		return f.(gin.HandlerFunc)
	} else {
		return gin.HandlerFunc(h)
	}
}

// @title           Hexagonal Saga API
// @version         0.1
// @description     hexagonal saga orchestration example api server

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGinRouter(cfg *config.Config) *GinRouter {
	logger := util.GetLogger().With(
		"module", "GinRouter",
	)

	if cfg.App.Stage != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.LeakBucket(&cfg.App))
	r.Use(middleware.Timeout(&cfg.App))
	r.Use(middleware.Cors())
	r.Use(middleware.XrayTracing([]string{"/healthy", "/swagger"}))
	r.Use(util.GinzapWithConfig(logger, &util.Config{
		UTC:       false,
		SkipPaths: []string{"/healthy"},
	}))
	r.Use(util.RecoveryWithZap(logger, true))

	r.GET("/healthy", func(c *gin.Context) {
		c.String(http.StatusOK, "OK\n")
	})

	if cfg.App.Stage == "local" || cfg.App.Stage == "dev" {
		docs.SwaggerInfo.BasePath = "/v1"
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return &GinRouter{
		r,
	}
}

func NewServer(cfg *config.Config, h http.Handler) *http.Server {
	if cfg.App.DisableHTTP {
		return nil
	}

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           h,
		ReadHeaderTimeout: 30 * time.Second,
	}
}

// NewServerForce create http.Server regardless of config
func NewServerForce(cfg *config.Config, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           h,
		ReadHeaderTimeout: 30 * time.Second,
	}
}
