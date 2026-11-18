package teacherdto

import (
	"time"

	"github.com/geniussociety/backend/internal/model"
)

// ProjectCreatePayload represents JSON data submitted when creating a project.
type ProjectCreatePayload struct {
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Instructions  string     `json:"instructions"`
	SubjectID     string     `json:"subjectId"`
	ClassIDs      []string   `json:"classIds"`
	Deadline      time.Time  `json:"deadline"`
	AllowResubmit bool       `json:"allowResubmit"`
}

// ProjectUpdatePayload represents JSON data submitted when updating a project.
type ProjectUpdatePayload struct {
	Title             string     `json:"title,omitempty"`
	Description       string     `json:"description,omitempty"`
	Instructions      string     `json:"instructions,omitempty"`
	ClassIDs          []string   `json:"classIds,omitempty"`
	Deadline          *time.Time `json:"deadline,omitempty"`
	AllowResubmit     *bool      `json:"allowResubmit,omitempty"`
	KeepExistingFiles *bool      `json:"keepExistingFiles,omitempty"`
}

// ProjectResponseDTO represents the detailed response format for a project.
type ProjectResponseDTO struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Instructions   string          `json:"instructions"`
	SubjectID      string          `json:"subjectId"`
	SubjectName    string          `json:"subjectName"`
	ClassIDs       []string        `json:"classIds"`
	ClassNames     []string        `json:"classNames"`
	Deadline       time.Time       `json:"deadline"`
	Attachments    []model.FileRef `json:"attachments"`
	AllowResubmit  bool            `json:"allowResubmit"`
	TotalStudents  int64           `json:"totalStudents"`
	SubmittedCount int64           `json:"submittedCount"`
	GradedCount    int64           `json:"gradedCount"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// ProjectListDTO represents an item in the projects list.
type ProjectListDTO struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	SubjectID      string    `json:"subjectId"`
	SubjectName    string    `json:"subjectName"`
	ClassIDs       []string  `json:"classIds"`
	ClassNames     []string  `json:"classNames"`
	Deadline       time.Time `json:"deadline"`
	TotalStudents  int64     `json:"totalStudents"`
	SubmittedCount int64     `json:"submittedCount"`
	GradedCount    int64     `json:"gradedCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

// ProjectMeta holds pagination metadata for projects.
type ProjectMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

// GradeProjectSubmissionPayload represents payload when grading a project submission.
type GradeProjectSubmissionPayload struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback"`
}

// ProjectSubmissionDTO represents a student's project submission item.
type ProjectSubmissionDTO struct {
	ID          string          `json:"id"`
	ProjectID   string          `json:"projectId"`
	StudentID   string          `json:"studentId"`
	StudentName string          `json:"studentName"`
	NIS         string          `json:"nis"`
	ClassID     string          `json:"classId"`
	ClassName   string          `json:"className"`
	Files       []model.FileRef `json:"files"`
	SubmittedAt *time.Time      `json:"submittedAt,omitempty"`
	Version     int             `json:"version"`
	Score       *float64        `json:"score,omitempty"`
	Feedback    string          `json:"feedback,omitempty"`
	GradedAt    *time.Time      `json:"gradedAt,omitempty"`
	Status      string          `json:"status"` // graded, ungraded, missing, late
	IsLate      bool            `json:"isLate"`
}

// ProjectSubmissionsMeta holds metadata for project submissions list.
type ProjectSubmissionsMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}
