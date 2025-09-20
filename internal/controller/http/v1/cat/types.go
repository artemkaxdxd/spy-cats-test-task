package cat

import (
	"backend/config"
	request "backend/internal/controller/http/request/cat"
	response "backend/internal/controller/http/response/cat"
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type (
	service interface {
		GetCats(ctx context.Context, breed string) ([]response.Cat, config.ServiceCode, error)
		GetCatByID(ctx context.Context, catID uint) (response.Cat, config.ServiceCode, error)
		CreateCat(ctx context.Context, body request.Cat) (response.Cat, config.ServiceCode, error)
		UpdateCat(ctx context.Context, body request.UpdateCat, catID uint) (config.ServiceCode, error)
		DeleteCat(ctx context.Context, catID uint) (config.ServiceCode, error)
	}

	validator interface {
		ValidateStruct(i any) (bool, error)
	}

	handler struct {
		svc       service
		l         *slog.Logger
		validator validator
	}
)

func InitHandler(
	g *gin.Engine,
	l *slog.Logger,
	svc service,
	validator validator,
) {
	h := handler{svc, l, validator}

	cats := g.Group("cats")
	{
		cats.GET("", h.getCats)
		cats.GET("/:cat_id", h.getCatByID)

		cats.POST("", h.createCat)
		cats.PATCH("/:cat_id", h.updateCat)
		cats.DELETE("/:cat_id", h.deleteCat)
	}
}
