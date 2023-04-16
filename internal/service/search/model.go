package search

import (
	"time"

	"github.com/artemmarkaryan/exlex-backend/internal/schema"
	"github.com/artemmarkaryan/exlex-backend/pkg/databaseutil"
	"github.com/google/uuid"
)

type CreateSearchRequest struct {
	Creator                uuid.UUID
	Name                   string
	Description            string
	Price                  float64
	RequiredWorkExperience int
	Deadline               *time.Time
	RequiredSpecialities   []string
	RequiredEducation      []string
}

type Search struct {
	ID                     uuid.UUID
	Name                   string
	Description            string
	Price                  float64
	RequiredWorkExperience int
	Deadline               *time.Time
	CreatedAt              time.Time
	RequiredSpecialities   []string
	RequiredEducation      []string
}

func (s *Search) fillFromRaw(dbo schema.SearchFullDataRaw) (err error) {
	educations, err := databaseutil.BytesToStrings(dbo.Education)
	if err != nil {
		return err
	}

	specialities, err := databaseutil.BytesToStrings(dbo.Speciality)
	if err != nil {
		return err
	}

	*s = Search{
		ID:                     dbo.ID,
		Name:                   dbo.Name,
		Description:            dbo.Description,
		Price:                  dbo.Price,
		RequiredWorkExperience: dbo.RequiredWorkExperience,
		Deadline:               dbo.Deadline,
		CreatedAt:              dbo.CreatedAt,
		RequiredEducation:      educations,
		RequiredSpecialities:   specialities,
	}

	return nil
}

type SearchApplicationRequest struct {
	SearchID uuid.UUID
	UserID   uuid.UUID
	Comment  *string
}

type ListApplicantsRequest struct {
	SearchID uuid.UUID
	UserID   uuid.UUID
}

type Applicant struct {
	UserID     uuid.UUID
	Education  string
	FullName   string
	Experience int
	Speciality []string
}

type Application struct {
	Applicant
	ID        uuid.UUID
	CreatedAt time.Time
	Comment   *string
}

func (a *Application) FillFromRaw(dbo schema.SearchApplicationRaw) error {
	specialities, err := databaseutil.BytesToStrings(dbo.Specialities)
	if err != nil {
		return err
	}

	var comment *string
	if dbo.Comment.Valid {
		*comment = dbo.Comment.String
	}

	*a = Application{
		Applicant: Applicant{
			UserID:     dbo.UserID,
			Education:  dbo.Education,
			FullName:   dbo.FullName,
			Experience: dbo.Experience,
			Speciality: specialities,
		},
		ID:        dbo.ApplicationID,
		CreatedAt: dbo.CreatedAt,
		Comment:   comment,
	}
	return nil
}
