package dto

import "go.mongodb.org/mongo-driver/v2/bson"

// Alias existing DTO structs to match service expectations
type ClassCreateDTO = CreateClassReq
type ClassUpdateDTO = UpdateClassReq

// ClassListFilter defines filters for listing classes with pagination.
type ClassListFilter struct {
	Code           *string       `json:"code,omitempty"`
	Name           *string       `json:"name,omitempty"`
	Status         *string       `json:"status,omitempty"`
	EducationLevel *string       `json:"educationLevel,omitempty"`
	Department     *string       `json:"department,omitempty"`
	AcademicYearID *string       `json:"academicYearId,omitempty"`
	Page           int64         `json:"page,omitempty"`
	Limit          int64         `json:"limit,omitempty"`
	SchoolID       bson.ObjectID `json:"-"`
}

// PaginatedResult is a generic container for paginated responses.
type PaginatedResult[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}
