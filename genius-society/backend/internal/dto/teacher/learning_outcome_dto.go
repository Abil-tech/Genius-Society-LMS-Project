package teacherdto

import (
	"time"
)

// LearningOutcomeCreatePayload represents the request body to create a Learning Outcome (PP).
type LearningOutcomeCreatePayload struct {
	SubjectID   string   `json:"subjectId"`
	ClassIDs    []string `json:"classIds"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Order       int      `json:"order,omitempty"`
}

// LearningOutcomeUpdatePayload represents the request body to update a Learning Outcome.
type LearningOutcomeUpdatePayload struct {
	ClassIDs    []string `json:"classIds,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Order       *int     `json:"order,omitempty"`
}

// LearningOutcomeItemMappingPayload represents the request body to map an item (assignment/project/assessment) to a PP.
type LearningOutcomeItemMappingPayload struct {
	LearningOutcomeID *string `json:"learningOutcomeId"`
}

// LearningOutcomeResponseDTO represents the detailed response format for a Learning Outcome.
type LearningOutcomeResponseDTO struct {
	ID          string    `json:"id"`
	SubjectID   string    `json:"subjectId"`
	SubjectName string    `json:"subjectName,omitempty"`
	ClassIDs    []string  `json:"classIds"`
	ClassNames  []string  `json:"classNames,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Order       int       `json:"order"`
	ItemCount   int64     `json:"itemCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
