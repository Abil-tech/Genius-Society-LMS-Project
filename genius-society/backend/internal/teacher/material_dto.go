package teacher

import (
	"time"

	"github.com/geniussociety/backend/internal/model"
)

// MaterialCreatePayload represents JSON data submitted when creating a material.
type MaterialCreatePayload struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	SubjectID   string     `json:"subjectId"`
	ClassIDs    []string   `json:"classIds"`
	Link        string     `json:"link,omitempty"`
	PublishNow  bool       `json:"publishNow"`
	PublishAt   *time.Time `json:"publishAt,omitempty"`
}

// MaterialUpdatePayload represents JSON data submitted when updating a material.
type MaterialUpdatePayload struct {
	Title            string     `json:"title,omitempty"`
	Description      string     `json:"description,omitempty"`
	ClassIDs         []string   `json:"classIds,omitempty"`
	Link             string     `json:"link,omitempty"`
	PublishNow       *bool      `json:"publishNow,omitempty"`
	PublishAt        *time.Time `json:"publishAt,omitempty"`
	KeepExistingFile *bool      `json:"keepExistingFile,omitempty"`
}

// MaterialResponseDTO represents the detailed response format for a material.
type MaterialResponseDTO struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	SubjectID   string          `json:"subjectId"`
	SubjectName string          `json:"subjectName"`
	ClassIDs    []string        `json:"classIds"`
	ClassNames  []string        `json:"classNames"`
	File        *model.FileRef  `json:"file,omitempty"`
	Attachments []model.FileRef `json:"attachments,omitempty"`
	Link        string          `json:"link,omitempty"`
	PublishAt   *time.Time      `json:"publishAt,omitempty"`
	IsPublished bool            `json:"isPublished"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// MaterialListDTO represents an item in the materials list.
type MaterialListDTO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	SubjectID   string     `json:"subjectId"`
	SubjectName string     `json:"subjectName"`
	ClassIDs    []string   `json:"classIds"`
	ClassNames  []string   `json:"classNames"`
	HasFile     bool       `json:"hasFile"`
	HasLink     bool       `json:"hasLink"`
	PublishAt   *time.Time `json:"publishAt,omitempty"`
	IsPublished bool       `json:"isPublished"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// MaterialMeta holds pagination metadata for materials.
type MaterialMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}
