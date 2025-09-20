package cat

import (
	"time"
)

// Could be split into different packages, but service is rather small,
// so to not overcomplicate leave it in one package for convenience.
type (
	Cat struct {
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `gorm:"index"`

		Name            string
		YearsExperience uint8
		Breed           string
		Salary          uint64 // Salary in cents (e.g. 100 = 1$ in cents)
	}

	Mission struct {
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `gorm:"index"`

		CatID       *uint `gorm:"index"`
		IsCompleted bool

		Cat     Cat      `gorm:"-"`
		Targets []Target `gorm:"-"`
	}

	Target struct {
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `gorm:"index"`

		MissionID   uint `gorm:"index"`
		Name        string
		Country     string
		Notes       string
		IsCompleted bool
	}
)

func (Cat) TableName() string {
	return "cats"
}

func (Mission) TableName() string {
	return "missions" // Could also be cat_missions, depending on the needed architecture
}

func (Target) TableName() string {
	return "targets" // Could also be (cat_)mission_targets, depending on the needed architecture
}
