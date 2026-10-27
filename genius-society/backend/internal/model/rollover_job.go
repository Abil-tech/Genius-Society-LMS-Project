package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobStatus string

const (
	JobStatusPending            JobStatus = "pending"
	JobStatusProcessing         JobStatus = "processing"
	JobStatusCompleted          JobStatus = "completed"
	JobStatusCompletedWithErrors JobStatus = "completed_with_errors"
	JobStatusFailed             JobStatus = "failed"
)

type RolloverJobSummary struct {
	TotalStudents int64 `bson:"totalStudents" json:"totalStudents"`
	Promoted      int64 `bson:"promoted" json:"promoted"`
	Retained      int64 `bson:"retained" json:"retained"`
	Graduated     int64 `bson:"graduated" json:"graduated"`
	Suspended     int64 `bson:"suspended" json:"suspended"`
	Excluded      int64 `bson:"excluded" json:"excluded"`
	Failed        int64 `bson:"failed" json:"failed"`
}

type RolloverJob struct {
	ID                   bson.ObjectID      `bson:"_id,omitempty" json:"id"`
	SchoolID             bson.ObjectID      `bson:"schoolId" json:"schoolId"`
	SourceAcademicYearID bson.ObjectID      `bson:"sourceAcademicYearId" json:"sourceAcademicYearId"`
	TargetAcademicYearID *bson.ObjectID     `bson:"targetAcademicYearId,omitempty" json:"targetAcademicYearId,omitempty"`
	ProcessType          ProcessType        `bson:"processType" json:"processType"` // promotion / graduation
	Status               JobStatus          `bson:"status" json:"status"`
	Summary              RolloverJobSummary `bson:"summary" json:"summary"`
	ExecutedBy           bson.ObjectID      `bson:"executedBy" json:"executedBy"`
	IdempotencyKey       string             `bson:"idempotencyKey" json:"idempotencyKey"`
	ErrorMessage         string             `bson:"errorMessage,omitempty" json:"errorMessage,omitempty"`
	StartedAt            time.Time          `bson:"startedAt" json:"startedAt"`
	CompletedAt          *time.Time         `bson:"completedAt,omitempty" json:"completedAt,omitempty"`
	CreatedAt            time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time          `bson:"updatedAt" json:"updatedAt"`
}
