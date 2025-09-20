package mission

import (
	entity "backend/internal/entity/cat"
	"backend/pkg/postgres"
	"context"
	"database/sql"
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

func (r repo) GetMissions(ctx context.Context) ([]entity.Mission, error) {
	rows, err := r.db.Instance().WithContext(ctx).Raw(`
		SELECT 
			m.id, m.created_at, m.updated_at, m.deleted_at,
			m.cat_id, m.is_completed,
			c.id, c.created_at, c.updated_at, c.deleted_at,
			c.name, c.years_experience, c.breed, c.salary,
			t.id, t.created_at, t.updated_at, t.deleted_at,
			t.mission_id, t.name, t.country, t.notes, t.is_completed
		FROM missions m
		LEFT JOIN cats c ON m.cat_id = c.id AND c.deleted_at IS NULL
		LEFT JOIN targets t ON m.id = t.mission_id AND t.deleted_at IS NULL
		WHERE m.deleted_at IS NULL
		ORDER BY m.created_at DESC, t.id ASC`).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	missionMap := make(map[uint]*entity.Mission)
	var orderedIDs []uint

	for rows.Next() {
		var (
			mission entity.Mission
			target  entity.Target

			catID, catYearsExperience, catSalary,
			targetID, targetMissionID sql.NullInt64

			catCreatedAt, catUpdatedAt,
			catDeletedAt, targetCreatedAt, targetUpdatedAt,
			targetDeletedAt sql.NullTime

			catName, catBreed, targetName,
			targetCountry, targetNotes sql.NullString

			targetIsCompleted sql.NullBool
		)

		if err = rows.Scan(
			&mission.ID, &mission.CreatedAt, &mission.UpdatedAt, &mission.DeletedAt,
			&mission.CatID, &mission.IsCompleted,
			&catID, &catCreatedAt, &catUpdatedAt, &catDeletedAt,
			&catName, &catYearsExperience, &catBreed, &catSalary,
			&targetID, &targetCreatedAt, &targetUpdatedAt, &targetDeletedAt,
			&targetMissionID, &targetName, &targetCountry, &targetNotes, &targetIsCompleted,
		); err != nil {
			return nil, err
		}

		v, ok := missionMap[mission.ID]
		if !ok {
			mission.Targets = make([]entity.Target, 0)

			if catID.Valid {
				mission.Cat = entity.Cat{
					ID:              uint(catID.Int64),
					CreatedAt:       catCreatedAt.Time,
					UpdatedAt:       catUpdatedAt.Time,
					Name:            catName.String,
					YearsExperience: uint8(catYearsExperience.Int64),
					Breed:           catBreed.String,
					Salary:          uint64(catSalary.Int64),
				}
				if catDeletedAt.Valid {
					deletedAt := catDeletedAt.Time
					mission.Cat.DeletedAt = &deletedAt
				}
			}

			missionMap[mission.ID] = &mission
			orderedIDs = append(orderedIDs, mission.ID)
			v = &mission
		}

		if targetID.Valid {
			target = entity.Target{
				ID:          uint(targetID.Int64),
				CreatedAt:   targetCreatedAt.Time,
				UpdatedAt:   targetUpdatedAt.Time,
				MissionID:   uint(targetMissionID.Int64),
				Name:        targetName.String,
				Country:     targetCountry.String,
				Notes:       targetNotes.String,
				IsCompleted: targetIsCompleted.Bool,
			}
			if targetDeletedAt.Valid {
				deletedAt := targetDeletedAt.Time
				target.DeletedAt = &deletedAt
			}
			v.Targets = append(v.Targets, target)
		}
	}

	missions := make([]entity.Mission, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		missions = append(missions, *missionMap[id])
	}

	return missions, nil
}

func (r repo) GetMission(ctx context.Context, missionID uint) (entity.Mission, error) {
	rows, err := r.db.Instance().WithContext(ctx).Raw(`
		SELECT 
			m.id, m.created_at, m.updated_at, m.deleted_at,
			m.cat_id, m.is_completed,
			c.id, c.created_at, c.updated_at, c.deleted_at,
			c.name, c.years_experience, c.breed, c.salary,
			t.id, t.created_at, t.updated_at, t.deleted_at,
			t.mission_id, t.name, t.country, t.notes, t.is_completed
		FROM missions m
		LEFT JOIN cats c ON m.cat_id = c.id AND c.deleted_at IS NULL
		LEFT JOIN targets t ON m.id = t.mission_id AND t.deleted_at IS NULL
		WHERE m.id = ? AND m.deleted_at IS NULL
		ORDER BY t.id ASC`,
		missionID).Rows()
	if err != nil {
		return entity.Mission{}, err
	}
	defer rows.Close()

	var mission entity.Mission
	mission.Targets = make([]entity.Target, 0)
	missionSet := false

	for rows.Next() {
		var (
			tempMission entity.Mission
			target      entity.Target

			catID, catYearsExperience, catSalary,
			targetID, targetMissionID sql.NullInt64

			catCreatedAt, catUpdatedAt,
			catDeletedAt, targetCreatedAt, targetUpdatedAt,
			targetDeletedAt sql.NullTime

			catName, catBreed, targetName,
			targetCountry, targetNotes sql.NullString

			targetIsCompleted sql.NullBool
		)

		if err = rows.Scan(
			&tempMission.ID, &tempMission.CreatedAt, &tempMission.UpdatedAt, &tempMission.DeletedAt,
			&tempMission.CatID, &tempMission.IsCompleted,
			&catID, &catCreatedAt, &catUpdatedAt, &catDeletedAt,
			&catName, &catYearsExperience, &catBreed, &catSalary,
			&targetID, &targetCreatedAt, &targetUpdatedAt, &targetDeletedAt,
			&targetMissionID, &targetName, &targetCountry, &targetNotes, &targetIsCompleted,
		); err != nil {
			return entity.Mission{}, err
		}

		if !missionSet {
			mission = tempMission
			mission.Targets = make([]entity.Target, 0)

			if catID.Valid {
				mission.Cat = entity.Cat{
					ID:              uint(catID.Int64),
					CreatedAt:       catCreatedAt.Time,
					UpdatedAt:       catUpdatedAt.Time,
					Name:            catName.String,
					YearsExperience: uint8(catYearsExperience.Int64),
					Breed:           catBreed.String,
					Salary:          uint64(catSalary.Int64),
				}
				if catDeletedAt.Valid {
					deletedAt := catDeletedAt.Time
					mission.Cat.DeletedAt = &deletedAt
				}
			}

			missionSet = true
		}

		if targetID.Valid {
			target = entity.Target{
				ID:          uint(targetID.Int64),
				CreatedAt:   targetCreatedAt.Time,
				UpdatedAt:   targetUpdatedAt.Time,
				MissionID:   uint(targetMissionID.Int64),
				Name:        targetName.String,
				Country:     targetCountry.String,
				Notes:       targetNotes.String,
				IsCompleted: targetIsCompleted.Bool,
			}
			if targetDeletedAt.Valid {
				deletedAt := targetDeletedAt.Time
				target.DeletedAt = &deletedAt
			}
			mission.Targets = append(mission.Targets, target)
		}
	}

	return mission, nil
}

func (r repo) CreateMission(ctx context.Context, tx *gorm.DB, mission entity.Mission) (entity.Mission, error) {
	err := tx.WithContext(ctx).Create(&mission).Error
	return mission, err
}

func (r repo) AssignCat(ctx context.Context, missionID, catID uint) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE missions
		SET cat_id = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		catID, time.Now(), missionID).Error
}

func (r repo) CompleteMission(ctx context.Context, missionID uint) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE missions
		SET is_completed = true, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		time.Now(), missionID).Error
}

func (r repo) DeleteMission(ctx context.Context, missionID uint) error {
	return r.db.Instance().WithContext(ctx).Exec(`
		UPDATE missions
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		time.Now(), missionID).Error
}
