package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProcessType string

const (
	ProcessTypePromotion  ProcessType = "promotion"
	ProcessTypeGraduation ProcessType = "graduation"
)

type EligibilityStatus string

const (
	EligibilityEligible      EligibilityStatus = "eligible"
	EligibilityPendingReview EligibilityStatus = "pending_review"
)

type DecisionStatus string

const (
	DecisionStatusPromoted  DecisionStatus = "promoted"
	DecisionStatusRetained  DecisionStatus = "retained"
	DecisionStatusGraduated DecisionStatus = "graduated"
	DecisionStatusSuspended DecisionStatus = "suspended"
	DecisionStatusExcluded  DecisionStatus = "excluded"
)

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusReviewed ReviewStatus = "reviewed"
	ReviewStatusResolved ReviewStatus = "resolved"
)

type AcademicDecision struct {
	ID                   bson.ObjectID     `bson:"_id,omitempty" json:"id"`
	SchoolID             bson.ObjectID     `bson:"schoolId" json:"schoolId"`
	StudentID            bson.ObjectID     `bson:"studentId" json:"studentId"`
	SourceAcademicYearID bson.ObjectID     `bson:"sourceAcademicYearId" json:"sourceAcademicYearId"`
	TargetAcademicYearID *bson.ObjectID    `bson:"targetAcademicYearId,omitempty" json:"targetAcademicYearId,omitempty"`
	SourceClassID        *bson.ObjectID    `bson:"sourceClassId,omitempty" json:"sourceClassId,omitempty"`
	TargetClassID        *bson.ObjectID    `bson:"targetClassId,omitempty" json:"targetClassId,omitempty"`
	ProcessType          ProcessType       `bson:"processType" json:"processType"`
	EligibilityStatus    EligibilityStatus `bson:"eligibilityStatus" json:"eligibilityStatus"`
	DecisionStatus       DecisionStatus    `bson:"decisionStatus" json:"decisionStatus"`
	ReasonCode           string            `bson:"reasonCode,omitempty" json:"reasonCode,omitempty"`
	Reason               string            `bson:"reason,omitempty" json:"reason,omitempty"`
	ReviewStatus         ReviewStatus      `bson:"reviewStatus" json:"reviewStatus"`
	ReviewedBy           *bson.ObjectID    `bson:"reviewedBy,omitempty" json:"reviewedBy,omitempty"`
	ReviewedAt           *time.Time        `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"`
	DecidedBy            bson.ObjectID     `bson:"decidedBy" json:"decidedBy"`
	DecidedAt            time.Time         `bson:"decidedAt" json:"decidedAt"`
	JobID                *bson.ObjectID    `bson:"jobId,omitempty" json:"jobId,omitempty"`
	CreatedAt            time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time         `bson:"updatedAt" json:"updatedAt"`
}
