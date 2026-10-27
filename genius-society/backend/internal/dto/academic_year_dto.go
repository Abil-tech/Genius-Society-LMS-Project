package dto

import "time"

type AcademicYearResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Semester  int       `json:"semester"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Create – used by Admin when adding a new year
type CreateAcademicYearDTO struct {
	Name      string    `json:"name" binding:"required"`
	Semester  int       `json:"semester" binding:"required,oneof=1 2"`
	StartDate time.Time `json:"startDate" binding:"required"`
	EndDate   time.Time `json:"endDate" binding:"required,gtfield=StartDate"`
}

// Update – partial update (PUT replaces all, PATCH can be a subset)
type UpdateAcademicYearDTO struct {
	Name      *string    `json:"name,omitempty" binding:"omitempty"`
	Semester  *int       `json:"semester,omitempty" binding:"omitempty,oneof=1 2"`
	StartDate *time.Time `json:"startDate,omitempty" binding:"omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty" binding:"omitempty,gtfield=StartDate"`
}
