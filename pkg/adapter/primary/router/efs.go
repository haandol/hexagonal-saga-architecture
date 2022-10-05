package router

import (
	"github.com/gin-gonic/gin"
	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/port/primaryport/routerport"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util/cerrors"
)

type EfsRouter struct {
	BaseRouter
	efsService *service.EfsService
}

func NewEfsRouter(
	efsService *service.EfsService,
) *EfsRouter {
	return &EfsRouter{
		efsService: efsService,
	}
}

func (r *EfsRouter) Route(rg routerport.RouterGroup) {
	g := rg.Group("/efs")
	g.Handle("GET", "/", r.WrappedHandler(r.ListHandler))
}

// @Summary list files
// @Schemes
// @Description list all files in the given path
// @Tags efs
// @Accept json
// @Produce json
// @Param path query string false "path"
// @Success 200 {array} string
// @Router /efs [get]
func (r *EfsRouter) ListHandler(c *gin.Context) *cerrors.CodedError {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	files, err := r.efsService.List(c.Request.Context(), path)
	if err != nil {
		return cerrors.New(constant.ErrFailToListFiles, err)
	}

	return r.Success(c, files)
}
