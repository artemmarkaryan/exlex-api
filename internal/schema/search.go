package schema

import (
	"time"

	"github.com/google/uuid"
)

type Search struct {
	ID                     uuid.UUID  `db:"id"`
	Creator                uuid.UUID  `db:"creator"`
	Name                   string     `db:"name"`
	Description            string     `db:"description"`
	Price                  float64    `db:"price"`
	RequiredWorkExperience int        `db:"required_work_experience"`
	Deadline               *time.Time `db:"deadline"`
	CreatedAt              time.Time  `db:"created_at"`
}

func (Search) TableName() string {
	return "search"
}

type SearchRequirementSpeciality struct {
	SearchUUID uuid.UUID `db:"search_uuid"`
	Speciality string    `db:"speciality"`
}

func (SearchRequirementSpeciality) TableName() string {
	return "search_requirement_speciality"
}

type SearchRequirementEducation struct {
	SearchUUID uuid.UUID `db:"search_uuid"`
	Education  string    `db:"education"`
}

func (SearchRequirementEducation) TableName() string {
	return "search_requirement_education"
}

type SearchFullDataRaw struct {
	ID                     uuid.UUID  `db:"id"`
	Creator                uuid.UUID  `db:"creator"`
	Name                   string     `db:"name"`
	Description            string     `db:"description"`
	Price                  float64    `db:"price"`
	RequiredWorkExperience int        `db:"required_work_experience"`
	Deadline               *time.Time `db:"deadline"`
	CreatedAt              time.Time  `db:"created_at"`
	Speciality             []byte     `db:"speciality"`
	Education              []byte     `db:"education"`
}

type SearchFullData struct {
	ID                     uuid.UUID
	Name                   string
	Description            string
	Price                  float64
	RequiredWorkExperience int
	Deadline               *time.Time
	CreatedAt              time.Time
	Speciality             []string
	Education              []string
}
