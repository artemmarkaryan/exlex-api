package schema

import (
	"github.com/google/uuid"
)

type Speciality struct {
	ID    string `db:"id"`
	Title string `db:"title"`
}

func (Speciality) TableName() string { return "speciality" }

type EducationType struct {
	ID    string `db:"id"`
	Title string `db:"title"`
}

func (EducationType) TableName() string { return "education_type" }

type UserSpeciality struct {
	UserUUID   uuid.UUID `db:"user_uuid"`
	Speciality string    `db:"speciality"`
}

func (UserSpeciality) TableName() string { return "user_speciality" }

type CustomerMetadata struct {
	UserUUID uuid.UUID `db:"user_uuid"`
	Name     string    `db:"name"`
}

func (CustomerMetadata) TableName() string { return "customer_metadata" }

type ExecutorMetadata struct {
	UserUUID        uuid.UUID `db:"user_uuid"`
	FullName        string    `db:"full_name"`
	ExperienceYears int       `db:"experience_years"`
	Education       string    `db:"education"`
}

func (ExecutorMetadata) TableName() string { return "executor_metadata" }

type FullExecutorMetadata struct {
	ExecutorMetadata
	Specialities []string
}
