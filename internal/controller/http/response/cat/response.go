package cat

import (
	entity "backend/internal/entity/cat"
	"time"
)

type (
	Cat struct {
		ID              uint   `json:"id"`
		Name            string `json:"name"`
		YearsExperience uint8  `json:"years_experience"`
		Breed           string `json:"breed"`
		Salary          uint64 `json:"salary"`

		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}

	Mission struct {
		ID          uint  `json:"id"`
		CatID       *uint `json:"cat_id"`
		IsCompleted bool  `json:"is_completed"`

		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`

		Cat     Cat      `json:"cat"`
		Targets []Target `json:"targets"`
	}

	Target struct {
		ID          uint   `json:"id"`
		MissionID   uint   `json:"mission_id"`
		Name        string `json:"name"`
		Country     string `json:"country"`
		Notes       string `json:"notes"`
		IsCompleted bool   `json:"is_completed"`

		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}
)

func CatToResponse(c entity.Cat) Cat {
	return Cat{
		ID:              c.ID,
		Name:            c.Name,
		YearsExperience: c.YearsExperience,
		Breed:           c.Breed,
		Salary:          c.Salary,

		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}

func CatsToResponse(cats []entity.Cat) []Cat {
	res := make([]Cat, 0, len(cats))
	for _, c := range cats {
		res = append(res, CatToResponse(c))
	}
	return res
}

func MissionToResponse(m entity.Mission) Mission {
	return Mission{
		ID:          m.ID,
		CatID:       m.CatID,
		IsCompleted: m.IsCompleted,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,

		Cat:     CatToResponse(m.Cat),
		Targets: TargetsToResponse(m.Targets),
	}
}

func MissionsToResponse(missions []entity.Mission) []Mission {
	res := make([]Mission, 0, len(missions))
	for _, m := range missions {
		res = append(res, MissionToResponse(m))
	}
	return res
}

func TargetToResponse(t entity.Target) Target {
	return Target{
		ID:          t.ID,
		MissionID:   t.MissionID,
		Name:        t.Name,
		Country:     t.Country,
		Notes:       t.Notes,
		IsCompleted: t.IsCompleted,

		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: t.DeletedAt,
	}
}

func TargetsToResponse(targets []entity.Target) []Target {
	res := make([]Target, 0, len(targets))
	for _, t := range targets {
		res = append(res, TargetToResponse(t))
	}
	return res
}
