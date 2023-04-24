package user_profile

import "github.com/google/uuid"

type UpdateUserProfileData struct {
	UserUUID uuid.UUID
}

type UpdateExecutorProfileData struct {
	UpdateUserProfileData
	FullName        *string
	ExperienceYears *int
	Specialities    []string
	Education       *string
}

type UpdateCustomerProfileData struct {
	UpdateUserProfileData
	Name *string
}

type CustomerProfile struct {
	FullName string
}

type ExecutorProfile struct {
	FullName        string
	WorkExperience  int
	EducationTypeID string
	Specialization  []string
}
