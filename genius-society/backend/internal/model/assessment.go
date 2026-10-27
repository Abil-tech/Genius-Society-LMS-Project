package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AssessmentCategory string
type AssessmentStatus string
type QuestionType string
type AttemptStatus string

const (
	CategoryDaily AssessmentCategory = "daily"
	CategoryUTS   AssessmentCategory = "uts"
	CategoryUAS   AssessmentCategory = "uas"

	AssessmentDraft     AssessmentStatus = "draft"
	AssessmentPublished AssessmentStatus = "published"
	AssessmentClosed    AssessmentStatus = "closed"

	QuestionMCQ   QuestionType = "mcq"
	QuestionEssay QuestionType = "essay"

	AttemptInProgress AttemptStatus = "in_progress"
	AttemptSubmitted  AttemptStatus = "submitted"
	AttemptGraded     AttemptStatus = "graded"
)

type Assessment struct {
	ID              bson.ObjectID      `bson:"_id,omitempty" json:"id"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	DurationMinutes int                `bson:"durationMinutes" json:"durationMinutes"`
	SubjectID       bson.ObjectID      `bson:"subjectId" json:"subjectId"`
	ClassIDs        []bson.ObjectID    `bson:"classIds" json:"classIds"`
	StartAt         *time.Time         `bson:"startAt,omitempty" json:"startAt,omitempty"`
	EndAt           *time.Time         `bson:"endAt,omitempty" json:"endAt,omitempty"`
	ShowResult      bool               `bson:"showResult" json:"showResult"`
	Category        AssessmentCategory `bson:"category" json:"category"`
	TeacherID       bson.ObjectID      `bson:"teacherId" json:"teacherId"`
	Status          AssessmentStatus   `bson:"status" json:"status"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type AssessmentQuestion struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AssessmentID  bson.ObjectID `bson:"assessmentId" json:"assessmentId"`
	Type          QuestionType  `bson:"type" json:"type"`
	Text          string        `bson:"text" json:"text"`
	Options       []string      `bson:"options,omitempty" json:"options,omitempty"`             // for MCQ
	CorrectOption *int          `bson:"correctOption,omitempty" json:"correctOption,omitempty"` // 0-based index for MCQ
	Weight        float64       `bson:"weight" json:"weight"`
	Order         int           `bson:"order" json:"order"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type AssessmentAttempt struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AssessmentID bson.ObjectID `bson:"assessmentId" json:"assessmentId"`
	StudentID    bson.ObjectID `bson:"studentId" json:"studentId"`
	StartedAt    time.Time     `bson:"startedAt" json:"startedAt"`
	ExpiresAt    time.Time     `bson:"expiresAt" json:"expiresAt"`
	SubmittedAt  *time.Time    `bson:"submittedAt,omitempty" json:"submittedAt,omitempty"`
	Status       AttemptStatus `bson:"status" json:"status"`
	AutoScore    float64       `bson:"autoScore" json:"autoScore"`
	ManualScore  float64       `bson:"manualScore" json:"manualScore"`
	TotalScore   float64       `bson:"totalScore" json:"totalScore"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type AssessmentAnswer struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AttemptID  bson.ObjectID `bson:"attemptId" json:"attemptId"`
	QuestionID bson.ObjectID `bson:"questionId" json:"questionId"`
	Answer     string        `bson:"answer" json:"answer"`
	IsCorrect  *bool         `bson:"isCorrect,omitempty" json:"isCorrect,omitempty"`
	Score      float64       `bson:"score" json:"score"`
	Feedback   string        `bson:"feedback,omitempty" json:"feedback,omitempty"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time     `bson:"updatedAt" json:"updatedAt"`
}
