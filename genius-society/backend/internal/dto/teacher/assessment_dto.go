package teacherdto

import (
	"time"

	"github.com/geniussociety/backend/internal/model"
)

// AssessmentCreatePayload represents payload to create an assessment.
type AssessmentCreatePayload struct {
	Title            string                     `json:"title"`
	Description      string                     `json:"description"`
	DurationMinutes  int                        `json:"durationMinutes"`
	SubjectID        string                     `json:"subjectId"`
	ClassIDs         []string                   `json:"classIds"`
	StartAt          *time.Time                 `json:"startAt,omitempty"`
	EndAt            *time.Time                 `json:"endAt,omitempty"`
	Category         model.AssessmentCategory   `json:"category"` // daily, uts, uas
	MaxAttempts      int                        `json:"maxAttempts,omitempty"`
	ShowResultMode   string                     `json:"showResultMode,omitempty"` // immediate, after_grading, never
	ShuffleQuestions bool                       `json:"shuffleQuestions,omitempty"`
}

// AssessmentUpdatePayload represents payload to update an assessment.
type AssessmentUpdatePayload struct {
	Title            string     `json:"title,omitempty"`
	Description      string     `json:"description,omitempty"`
	DurationMinutes  *int       `json:"durationMinutes,omitempty"`
	ClassIDs         []string   `json:"classIds,omitempty"`
	StartAt          *time.Time `json:"startAt,omitempty"`
	EndAt            *time.Time `json:"endAt,omitempty"`
	MaxAttempts      *int       `json:"maxAttempts,omitempty"`
	ShowResultMode   string     `json:"showResultMode,omitempty"`
	ShuffleQuestions *bool      `json:"shuffleQuestions,omitempty"`
}

// QuestionCreateUpdatePayload represents payload to create or update a question.
type QuestionCreateUpdatePayload struct {
	Type          model.QuestionType `json:"type"` // mcq, essay
	Text          string             `json:"text"`
	Options       []string           `json:"options,omitempty"`
	CorrectOption *int               `json:"correctOption,omitempty"`
	Weight        float64            `json:"weight"`
	Order         int                `json:"order,omitempty"`
}

// QuestionReorderItem represents an item when reordering questions.
type QuestionReorderItem struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
}

// QuestionReorderPayload represents payload to reorder questions.
type QuestionReorderPayload struct {
	Questions []QuestionReorderItem `json:"questions"`
}

// RegradeMCQPayload represents payload to regrade an MCQ question after attempts exist.
type RegradeMCQPayload struct {
	QuestionID       string `json:"questionId"`
	NewCorrectOption int    `json:"newCorrectOption"`
}

// GradeEssayAnswerPayload represents payload when grading a single essay answer.
type GradeEssayAnswerPayload struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback"`
}

// QuestionTeacherDTO represents question detail visible to teachers (includes correct answer).
type QuestionTeacherDTO struct {
	ID            string             `json:"id"`
	AssessmentID  string             `json:"assessmentId"`
	Type          model.QuestionType `json:"type"`
	Text          string             `json:"text"`
	Options       []string           `json:"options,omitempty"`
	CorrectOption *int               `json:"correctOption,omitempty"`
	Weight        float64            `json:"weight"`
	Order         int                `json:"order"`
	CreatedAt     time.Time          `json:"createdAt"`
}

// AssessmentResponseDTO represents assessment detail for teachers.
type AssessmentResponseDTO struct {
	ID               string                   `json:"id"`
	Title            string                   `json:"title"`
	Description      string                   `json:"description"`
	DurationMinutes  int                      `json:"durationMinutes"`
	SubjectID        string                   `json:"subjectId"`
	SubjectName      string                   `json:"subjectName"`
	ClassIDs         []string                 `json:"classIds"`
	ClassNames       []string                 `json:"classNames"`
	StartAt          *time.Time               `json:"startAt,omitempty"`
	EndAt            *time.Time               `json:"endAt,omitempty"`
	Category         model.AssessmentCategory `json:"category"`
	Status           model.AssessmentStatus   `json:"status"`
	MaxAttempts      int                      `json:"maxAttempts"`
	ShowResultMode   string                   `json:"showResultMode"`
	ShuffleQuestions bool                     `json:"shuffleQuestions"`
	TotalScore       float64                  `json:"totalScore"`
	QuestionCount    int                      `json:"questionCount"`
	Questions        []QuestionTeacherDTO     `json:"questions,omitempty"`
	AttemptCount     int64                    `json:"attemptCount"`
	CreatedAt        time.Time                `json:"createdAt"`
	UpdatedAt        time.Time                `json:"updatedAt"`
}

// AssessmentListDTO represents assessment item in list view.
type AssessmentListDTO struct {
	ID              string                   `json:"id"`
	Title           string                   `json:"title"`
	SubjectID       string                   `json:"subjectId"`
	SubjectName     string                   `json:"subjectName"`
	ClassIDs        []string                 `json:"classIds"`
	ClassNames      []string                 `json:"classNames"`
	DurationMinutes int                      `json:"durationMinutes"`
	Category        model.AssessmentCategory `json:"category"`
	Status          model.AssessmentStatus   `json:"status"`
	StartAt         *time.Time               `json:"startAt,omitempty"`
	EndAt           *time.Time               `json:"endAt,omitempty"`
	QuestionCount   int                      `json:"questionCount"`
	TotalScore      float64                  `json:"totalScore"`
	AttemptCount    int64                    `json:"attemptCount"`
	CreatedAt       time.Time                `json:"createdAt"`
}

// AssessmentMeta holds pagination metadata for assessments.
type AssessmentMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

// AssessmentAttemptDTO represents student's attempt summary for teachers.
type AssessmentAttemptDTO struct {
	ID            string              `json:"id"`
	AssessmentID  string              `json:"assessmentId"`
	StudentID     string              `json:"studentId"`
	StudentName   string              `json:"studentName"`
	NIS           string              `json:"nis"`
	ClassID       string              `json:"classId"`
	ClassName     string              `json:"className"`
	AttemptNumber int                 `json:"attemptNumber"`
	Status        model.AttemptStatus `json:"status"` // in_progress, submitted, graded
	StartedAt     time.Time           `json:"startedAt"`
	ExpiresAt     time.Time           `json:"expiresAt"`
	SubmittedAt   *time.Time          `json:"submittedAt,omitempty"`
	AutoScore     float64             `json:"autoScore"`
	ManualScore   float64             `json:"manualScore"`
	TotalScore    float64             `json:"totalScore"`
	PendingEssays int                 `json:"pendingEssays"`
}

// AnswerDetailDTO represents answer detail inside attempt view.
type AnswerDetailDTO struct {
	ID            string             `json:"id"`
	QuestionID    string             `json:"questionId"`
	QuestionType  model.QuestionType `json:"questionType"`
	QuestionText  string             `json:"questionText"`
	Options       []string           `json:"options,omitempty"`
	CorrectOption *int               `json:"correctOption,omitempty"`
	Weight        float64            `json:"weight"`
	Answer        string             `json:"answer"`
	IsCorrect     *bool              `json:"isCorrect,omitempty"`
	Score         float64            `json:"score"`
	Feedback      string             `json:"feedback,omitempty"`
	Order         int                `json:"order"`
}

// AttemptDetailResponseDTO represents detailed attempt view for grading.
type AttemptDetailResponseDTO struct {
	Attempt AssessmentAttemptDTO `json:"attempt"`
	Answers []AnswerDetailDTO    `json:"answers"`
}

// AttemptsMeta holds pagination metadata for attempts.
type AttemptsMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}
