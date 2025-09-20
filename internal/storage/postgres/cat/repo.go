package cat

import (
	entity "backend/internal/entity/cat"
	"backend/pkg/postgres"
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type repo struct {
	db postgres.Database
}

func NewRepo(db postgres.Database) repo {
	return repo{db}
}

func (r repo) GetDB(ctx context.Context) *gorm.DB {
	return r.db.Instance().WithContext(ctx)
}

func (r repo) NewTransaction(ctx context.Context) *gorm.DB {
	return r.db.Instance().WithContext(ctx).Begin()
}

func (r repo) GetCats(ctx context.Context, breed string) (cats []entity.Cat, err error) {
	var (
		query strings.Builder
		args  []any
	)

	query.WriteString("SELECT * FROM cats WHERE deleted_at IS NULL")

	// Breed is a hardcoded filter, if more filters needed
	// pass map instead with needed field names and value to filter
	if breed != "" {
		query.WriteString(" AND breed = ?")
		args = append(args, breed)
	}

	err = r.db.Instance().WithContext(ctx).
		Raw(query.String(), args...).Scan(&cats).Error
	return
}

func (r repo) GetCatByID(ctx context.Context, catID uint) (cat entity.Cat, err error) {
	err = r.db.Instance().WithContext(ctx).Raw(`
		SELECT * FROM cats
		WHERE id = ? AND deleted_at IS NULL`,
		catID).Scan(&cat).Error
	return
}

func (r repo) CreateCat(ctx context.Context, cat entity.Cat) (entity.Cat, error) {
	err := r.db.Instance().WithContext(ctx).Create(&cat).Error
	return cat, err
}

func (r repo) UpdateCat(ctx context.Context, cat entity.Cat) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE cats
		SET salary = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		cat.Salary, time.Now(),
		cat.ID).Error
}

func (r repo) DeleteCat(ctx context.Context, catID uint) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE cats
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		time.Now(), catID).Error
}
