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

		GetCats(ctx context.Context, breed string) (cats []entity.Cat, err error)
		GetCatByID(ctx context.Context, catID uint) (cat entity.Cat, err error)

		CreateCat(ctx context.Context, cat entity.Cat) (entity.Cat, error)
		UpdateCat(ctx context.Context, cat entity.Cat) error
		DeleteCat(ctx context.Context, catID uint) error
	}

	breedValidator interface {
		IsValid(ctx context.Context, breedName string) (bool, error)
	}

	service struct {
		repo           repo
		breedValidator breedValidator
		l              *slog.Logger
	}
)

func NewService(
	repo repo,
	breedValidator breedValidator,
	l *slog.Logger,
) service {
	return service{
		repo,
		breedValidator,
		l}
}
