package target

import (
	entity "backend/internal/entity/cat"
	"backend/pkg/postgres"
	"context"
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

func (r repo) GetTargetByID(ctx context.Context, targetID, missionID uint) (target entity.Target, err error) {
	err = r.db.Instance().WithContext(ctx).Raw(`
		SELECT * FROM targets
		WHERE id = ? AND mission_id = ? AND deleted_at IS NULL`,
		targetID, missionID).Scan(&target).Error
	return
}

func (r repo) GetTargetsByMissionID(ctx context.Context, missionID uint) (targets []entity.Target, err error) {
	err = r.db.Instance().WithContext(ctx).Raw(`
		SELECT * FROM targets
		WHERE mission_id = ? AND deleted_at IS NULL`,
		missionID).Scan(&targets).Error
	return
}

func (r repo) CreateTarget(ctx context.Context, target entity.Target) error {
	return r.db.Instance().WithContext(ctx).Create(&target).Error
}

func (r repo) CreateTargets(ctx context.Context, tx *gorm.DB, targets []entity.Target) ([]entity.Target, error) {
	err := tx.WithContext(ctx).Create(&targets).Error
	return targets, err
}

func (r repo) UpdateTarget(ctx context.Context, target entity.Target) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE targets
		SET notes = ?, updated_at = ?
		WHERE id = ? AND mission_id = ? AND deleted_at IS NULL`,
		target.Notes, time.Now(),
		target.ID, target.MissionID).Error
}

func (r repo) CompleteTarget(ctx context.Context, targetID, missionID uint) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE targets
		SET is_completed = true, updated_at = ?
		WHERE id = ? AND mission_id = ? AND deleted_at IS NULL`,
		time.Now(), targetID, missionID).Error
}

func (r repo) DeleteTarget(ctx context.Context, targetID, missionID uint) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE targets
		SET deleted_at = ?
		WHERE id = ? AND mission_id = ? AND deleted_at IS NULL`,
		time.Now(), targetID, missionID).Error
}
