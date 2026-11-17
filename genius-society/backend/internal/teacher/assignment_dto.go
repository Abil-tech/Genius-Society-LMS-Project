package teacher

import (
	"time"

	"github.com/geniussociety/backend/internal/model"
)

// AssignmentCreatePayload represents JSON data submitted when creating an assignment.
type AssignmentCreatePayload struct {
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	SubjectID     string     `json:"subjectId"`
	ClassIDs      []string   `json:"classIds"`
	Deadline      time.Time  `json:"deadline"`
	AllowResubmit bool       `json:"allowResubmit"`
}

// AssignmentUpdatePayload represents JSON data submitted when updating an assignment.
type AssignmentUpdatePayload struct {
	Title            string     `json:"title,omitempty"`
	Description      string     `json:"description,omitempty"`
	ClassIDs         []string   `json:"classIds,omitempty"`
	Deadline         *time.Time `json:"deadline,omitempty"`
	AllowResubmit    *bool      `json:"allowResubmit,omitempty"`
	KeepExistingFiles *bool     `json:"keepExistingFiles,omitempty"`
}

// AssignmentResponseDTO represents the detailed response format for an assignment.
type AssignmentResponseDTO struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	SubjectID     string          `json:"subjectId"`
	SubjectName   string          `json:"subjectName"`
	ClassIDs      []string        `json:"classIds"`
	ClassNames    []string        `json:"classNames"`
	Deadline      time.Time       `json:"deadline"`
	Attachments   []model.FileRef `json:"attachments"`
	AllowResubmit bool            `json:"allowResubmit"`
	TotalStudents int64           `json:"totalStudents"`
	SubmittedCount int64          `json:"submittedCount"`
	GradedCount   int64           `json:"gradedCount"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

// AssignmentListDTO represents an item in the assignments list.
type AssignmentListDTO struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	SubjectID     string    `json:"subjectId"`
	SubjectName   string    `json:"subjectName"`
	ClassIDs      []string  `json:"classIds"`
	ClassNames    []string  `json:"classNames"`
	Deadline      time.Time `json:"deadline"`
	TotalStudents int64     `json:"totalStudents"`
	SubmittedCount int64    `json:"submittedCount"`
	GradedCount   int64     `json:"gradedCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

// AssignmentMeta holds pagination metadata for assignments.
type AssignmentMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

// GradeSubmissionPayload represents payload when teacher grades an assignment submission.
type GradeSubmissionPayload struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback"`
}

// AssignmentSubmissionDTO represents a student's submission item.
type AssignmentSubmissionDTO struct {
	ID           string          `json:"id"`
	AssignmentID string          `json:"assignmentId"`
	StudentID    string          `json:"studentId"`
	StudentName  string          `json:"studentName"`
	NIS          string          `json:"nis"`
	ClassID      string          `json:"classId"`
	ClassName    string          `json:"className"`
	Files        []model.FileRef `json:"files"`
	Note         string          `json:"note,omitempty"`
	SubmittedAt  *time.Time      `json:"submittedAt,omitempty"`
	Version      int             `json:"version"`
	Score        *float64        `json:"score,omitempty"`
	Feedback     string          `json:"feedback,omitempty"`
	GradedAt     *time.Time      `json:"gradedAt,omitempty"`
	Status       string          `json:"status"` // graded, ungraded, missing, late
	IsLate       bool            `json:"isLate"`
}

// SubmissionsMeta holds metadata for submissions list.
type SubmissionsMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}
