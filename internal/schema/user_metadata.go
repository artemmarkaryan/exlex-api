package schema

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MetadataKey string

const (
	MetadataKeyFullName       MetadataKey = "fullname"
	MetadataKeyWorkExperience MetadataKey = "work_experience"
)

type UserMetadata struct {
	UserUUID  uuid.UUID       `db:"user_uuid"`
	Metadata  json.RawMessage `db:"metadata"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (UserMetadata) TableName() string { return "user_metadata" }

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

type UserEducation struct {
	UserUUID  uuid.UUID `db:"user_uuid"`
	Education string    `db:"education"`
}

func (UserEducation) TableName() string { return "user_education" }
