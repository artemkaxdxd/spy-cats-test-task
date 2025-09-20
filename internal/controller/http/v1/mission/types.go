package mission

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
		GetMissions(ctx context.Context) ([]response.Mission, config.ServiceCode, error)
		GetMission(ctx context.Context, missionID uint) (response.Mission, config.ServiceCode, error)

		CreateMission(ctx context.Context, mission request.Mission) (response.Mission, config.ServiceCode, error)
		AssignCat(ctx context.Context, missionID, catID uint) (config.ServiceCode, error)
		CompleteMission(ctx context.Context, missionID uint) (config.ServiceCode, error)
		DeleteMission(ctx context.Context, missionID uint) (config.ServiceCode, error)

		CreateTarget(ctx context.Context, body request.Target, missionID uint) (config.ServiceCode, error)
		UpdateTarget(ctx context.Context, body request.UpdateTarget, targetID, missionID uint) (config.ServiceCode, error)
		CompleteTarget(ctx context.Context, targetID, missionID uint) (config.ServiceCode, error)
		DeleteTarget(ctx context.Context, targetID, missionID uint) (config.ServiceCode, error)
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

	missions := g.Group("missions")
	{
		missions.GET("", h.getMissions)
		missions.GET("/:mission_id", h.getMission)

		missions.POST("", h.createMission)

		missions.PATCH("/:mission_id/assign/:cat_id", h.assignCat)
		missions.PATCH("/:mission_id/complete", h.completeMission)

		missions.DELETE("/:mission_id", h.deleteMission)

		missions.POST("/:mission_id/targets", h.createTarget)
		missions.PATCH("/:mission_id/targets/:target_id", h.updateTarget)
		missions.PATCH("/:mission_id/targets/:target_id/complete", h.completeTarget)
		missions.DELETE("/:mission_id/targets/:target_id", h.deleteTarget)
	}
}
