package cat

import (
	entity "backend/internal/entity/cat"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type (
	repo interface {
		GetDB(ctx context.Context) *gorm.DB
		NewTransaction(ctx context.Context) *gorm.DB

		GetMissions(ctx context.Context) ([]entity.Mission, error)
		GetMission(ctx context.Context, missionID uint) (entity.Mission, error)

		CreateMission(ctx context.Context, tx *gorm.DB, mission entity.Mission) (entity.Mission, error)
		AssignCat(ctx context.Context, missionID, catID uint) error
		CompleteMission(ctx context.Context, missionID uint) error
		DeleteMission(ctx context.Context, missionID uint) error
	}

	targetRepo interface {
		GetTargetByID(ctx context.Context, targetID, missionID uint) (entity.Target, error)
		GetTargetsByMissionID(ctx context.Context, missionID uint) ([]entity.Target, error)

		CreateTarget(ctx context.Context, target entity.Target) error
		CreateTargets(ctx context.Context, tx *gorm.DB, targets []entity.Target) ([]entity.Target, error)
		UpdateTarget(ctx context.Context, target entity.Target) error
		CompleteTarget(ctx context.Context, targetID, missionID uint) error
		DeleteTarget(ctx context.Context, targetID, missionID uint) error
	}

	catRepo interface {
		GetCatByID(ctx context.Context, catID uint) (entity.Cat, error)
	}

	service struct {
		repo       repo
		targetRepo targetRepo
		catRepo    catRepo
		l          *slog.Logger
	}
)

func NewService(
	repo repo,
	targetRepo targetRepo,
	catRepo catRepo,
	l *slog.Logger,
) service {
	return service{repo, targetRepo, catRepo, l}
}
