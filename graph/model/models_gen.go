// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type CreateSearchInput struct {
	Title        string                   `json:"title"`
	Description  string                   `json:"description"`
	Price        float64                  `json:"price"`
	Deadline     *DateInput               `json:"deadline"`
	Requirements *SearchRequirementsInput `json:"requirements"`
}

type Customer struct {
	FullName string `json:"fullName"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type DateInput struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type EducationType struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Executor struct {
	FullName        *string  `json:"fullName"`
	WorkExperience  *int     `json:"workExperience"`
	EducationTypeID *string  `json:"educationTypeID"`
	Specialization  []string `json:"specialization"`
}

type Search struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	Price        float64             `json:"price"`
	CreatedAt    string              `json:"createdAt"`
	Deadline     *Date               `json:"deadline"`
	Requirements *SearchRequirements `json:"requirements"`
}

type SearchRequirements struct {
	Speciality     []string `json:"speciality"`
	EducationType  []string `json:"educationType"`
	WorkExperience int      `json:"workExperience"`
}

type SearchRequirementsInput struct {
	Speciality     []string `json:"speciality"`
	EducationType  []string `json:"educationType"`
	WorkExperience int      `json:"workExperience"`
}

type SetCustomerProfileData struct {
	FullName *string `json:"fullName"`
}

type SetExecutorProfileData struct {
	FullName        *string  `json:"fullName"`
	WorkExperience  *int     `json:"workExperience"`
	EducationTypeID *string  `json:"educationTypeID"`
	Specialization  []string `json:"specialization"`
}

type Speciality struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Role string

const (
	RoleExecutor Role = "executor"
	RoleCustomer Role = "customer"
)

var AllRole = []Role{
	RoleExecutor,
	RoleCustomer,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleExecutor, RoleCustomer:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
