package cat

import (
	"backend/config"
	entity "backend/internal/entity/cat"
	"fmt"
	"time"
)

type (
	Cat struct {
		Name            string `json:"name" valid:"required"`
		YearsExperience uint8  `json:"years_experience"`
		Breed           string `json:"breed" valid:"required"`
		Salary          uint64 `json:"salary"`
	}

	UpdateCat struct {
		Salary uint64 `json:"salary" valid:"required"`
	}

	Mission struct {
		CatID   *uint   `json:"cat_id"`
		Targets Targets `json:"targets"`
	}

	Target struct {
		Name    string `json:"name" binding:"required"`
		Country string `json:"country" binding:"required"`
		Notes   string `json:"notes"`
	}

	UpdateTarget struct {
		Notes string `json:"notes"`
	}

	Targets []Target
)

func (c Cat) ToEntity() entity.Cat {
	return entity.Cat{
		Name:            c.Name,
		YearsExperience: c.YearsExperience,
		Breed:           c.Breed,
		Salary:          c.Salary,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c UpdateCat) ToEntity(catID uint) entity.Cat {
	return entity.Cat{
		ID:     catID,
		Salary: c.Salary,
	}
}

func (m Mission) ToEntity() entity.Mission {
	return entity.Mission{
		CatID: m.CatID,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t Target) ToEntity(missionID uint) entity.Target {
	return entity.Target{
		MissionID: missionID,
		Name:      t.Name,
		Country:   t.Country,
		Notes:     t.Notes,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t Targets) ToEntity(missionID uint) []entity.Target {
	res := make([]entity.Target, 0, len(t))
	for _, v := range t {
		res = append(res, v.ToEntity(missionID))
	}
	return res
}

func (t UpdateTarget) ToEntity(targetID, missionID uint) entity.Target {
	return entity.Target{
		ID:        targetID,
		MissionID: missionID,
		Notes:     t.Notes,
	}
}

func (m Mission) ValidateTargetsLen() error {
	targetsLen := len(m.Targets)

	if targetsLen < config.MinMissionTargets ||
		targetsLen > config.MaxMissionTargets {
		return fmt.Errorf("mission targets amount should be in range (%d|%d): %d",
			config.MinMissionTargets, config.MaxMissionTargets, targetsLen)
	}

	return nil
}
