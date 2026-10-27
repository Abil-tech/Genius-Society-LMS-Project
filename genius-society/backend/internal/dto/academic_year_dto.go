package dto

import "time"

type AcademicYearResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Semester  int       `json:"semester"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	IsActive  bool      `json:"isActive"`
	Status    string    `json:"status"` // active, archived, inactive
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

// Promotion / Rollover DTOs
type PromotionCandidateItem struct {
	StudentID           string   `json:"studentId"`
	NIS                 string   `json:"nis"`
	Name                string   `json:"name"`
	SourceClassID       string   `json:"sourceClassId"`
	SourceClassName     string   `json:"sourceClassName"`
	EducationLevel      string   `json:"educationLevel"`
	EligibilityStatus   string   `json:"eligibilityStatus"` // eligible, pending_review
	RecommendedAction   string   `json:"recommendedAction"` // promoted, retained, suspended
	ReasonCode          string   `json:"reasonCode,omitempty"`
	Reason              string   `json:"reason,omitempty"`
	SuggestedClassID    string   `json:"suggestedClassId,omitempty"`
	TargetClassID       string   `json:"targetClassId,omitempty"`
	AverageScore        float64  `json:"averageScore"`
	MissingGradesCount  int      `json:"missingGradesCount"`
}

type PromotionDecisionInput struct {
	StudentID     string `json:"studentId" binding:"required"`
	Action        string `json:"action" binding:"required,oneof=promoted retained suspended excluded"`
	TargetClassID string `json:"targetClassId,omitempty"`
	ReasonCode    string `json:"reasonCode,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

type PromotionPreviewRequest struct {
	SourceAcademicYearID string                   `json:"sourceAcademicYearId" binding:"required"`
	TargetAcademicYearID string                   `json:"targetAcademicYearId" binding:"required"`
	Decisions            []PromotionDecisionInput `json:"decisions,omitempty"`
}

type PromotionExecutionRequest struct {
	SourceAcademicYearID string                   `json:"sourceAcademicYearId" binding:"required"`
	TargetAcademicYearID string                   `json:"targetAcademicYearId" binding:"required"`
	Decisions            []PromotionDecisionInput `json:"decisions" binding:"required,min=1,dive"`
	IdempotencyKey       string                   `json:"idempotencyKey" binding:"required"`
}

type RolloverJobSummaryDTO struct {
	TotalStudents int64 `json:"totalStudents"`
	Promoted      int64 `json:"promoted"`
	Retained      int64 `json:"retained"`
	Graduated     int64 `json:"graduated"`
	Suspended     int64 `json:"suspended"`
	Excluded      int64 `json:"excluded"`
	Failed        int64 `json:"failed"`
}

type ProcessPreviewResponse struct {
	TotalStudents int64                    `json:"totalStudents"`
	Eligible      int64                    `json:"eligible"`
	PendingReview int64                    `json:"pendingReview"`
	Promoted      int64                    `json:"promoted"`
	Retained      int64                    `json:"retained"`
	Graduated     int64                    `json:"graduated"`
	Suspended     int64                    `json:"suspended"`
	Excluded      int64                    `json:"excluded"`
	Errors        []string                 `json:"errors"`
	Candidates    []PromotionCandidateItem `json:"candidates,omitempty"`
}

type ProcessExecutionResponse struct {
	JobID       string                `json:"jobId"`
	Status      string                `json:"status"`
	Summary     RolloverJobSummaryDTO `json:"summary"`
	ExecutedAt  time.Time             `json:"executedAt"`
	Errors      []string              `json:"errors,omitempty"`
}

// Graduation DTOs
type GraduationCandidateItem struct {
	StudentID          string   `json:"studentId"`
	NIS                string   `json:"nis"`
	Name               string   `json:"name"`
	ClassID            string   `json:"classId"`
	ClassName          string   `json:"className"`
	EducationLevel     string   `json:"educationLevel"`
	EnrollmentYear     int      `json:"enrollmentYear"`
	EligibilityStatus  string   `json:"eligibilityStatus"` // eligible, pending_review
	RecommendedAction  string   `json:"recommendedAction"` // graduated, suspended
	ReasonCode         string   `json:"reasonCode,omitempty"`
	Reason             string   `json:"reason,omitempty"`
	AverageScore       float64  `json:"averageScore"`
	MissingGradesCount int      `json:"missingGradesCount"`
}

type GraduationDecisionInput struct {
	StudentID      string     `json:"studentId" binding:"required"`
	Action         string     `json:"action" binding:"required,oneof=graduated suspended excluded"`
	GraduationDate *time.Time `json:"graduationDate,omitempty"`
	ReasonCode     string     `json:"reasonCode,omitempty"`
	Reason         string     `json:"reason,omitempty"`
}

type GraduationPreviewRequest struct {
	AcademicYearID string                    `json:"academicYearId" binding:"required"`
	Decisions      []GraduationDecisionInput `json:"decisions,omitempty"`
}

type GraduationExecutionRequest struct {
	AcademicYearID string                    `json:"academicYearId" binding:"required"`
	Decisions      []GraduationDecisionInput `json:"decisions" binding:"required,min=1,dive"`
	IdempotencyKey string                    `json:"idempotencyKey" binding:"required"`
}

type GraduationPreviewResponse struct {
	TotalStudents int64                     `json:"totalStudents"`
	Eligible      int64                     `json:"eligible"`
	PendingReview int64                     `json:"pendingReview"`
	Graduated     int64                     `json:"graduated"`
	Suspended     int64                     `json:"suspended"`
	Excluded      int64                     `json:"excluded"`
	Errors        []string                  `json:"errors"`
	Candidates    []GraduationCandidateItem `json:"candidates,omitempty"`
}

// Suspended Decisions DTOs
type SuspendedDecisionItem struct {
	ID                   string     `json:"id"`
	StudentID            string     `json:"studentId"`
	StudentName          string     `json:"studentName"`
	NIS                  string     `json:"nis"`
	SourceClassName      string     `json:"sourceClassName"`
	SourceAcademicYearID string     `json:"sourceAcademicYearId"`
	TargetAcademicYearID string     `json:"targetAcademicYearId,omitempty"`
	ProcessType          string     `json:"processType"` // promotion / graduation
	EligibilityStatus    string     `json:"eligibilityStatus"`
	DecisionStatus       string     `json:"decisionStatus"`
	ReasonCode           string     `json:"reasonCode"`
	Reason               string     `json:"reason"`
	ReviewStatus         string     `json:"reviewStatus"`
	ReviewedByName       string     `json:"reviewedByName,omitempty"`
	ReviewedAt           *time.Time `json:"reviewedAt,omitempty"`
	DecidedAt            time.Time  `json:"decidedAt"`
	CreatedAt            time.Time  `json:"createdAt"`
}

type ReviewDecisionRequest struct {
	ReviewNotes string `json:"reviewNotes" binding:"required"`
}

type ResolveDecisionRequest struct {
	Action         string     `json:"action" binding:"required,oneof=promoted retained graduated excluded"`
	TargetClassID  string     `json:"targetClassId,omitempty"`  // required if promoted or retained
	GraduationDate *time.Time `json:"graduationDate,omitempty"` // optional for graduated
	Notes          string     `json:"notes,omitempty"`
}

type RolloverJobResponse struct {
	ID                     string                `json:"id"`
	ProcessType            string                `json:"processType"`
	SourceAcademicYearID   string                `json:"sourceAcademicYearId"`
	SourceAcademicYearName string                `json:"sourceAcademicYearName"`
	TargetAcademicYearID   string                `json:"targetAcademicYearId,omitempty"`
	TargetAcademicYearName string                `json:"targetAcademicYearName,omitempty"`
	Status                 string                `json:"status"`
	Summary                RolloverJobSummaryDTO `json:"summary"`
	ExecutedByName         string                `json:"executedByName"`
	ErrorMessage           string                `json:"errorMessage,omitempty"`
	StartedAt              time.Time             `json:"startedAt"`
	CompletedAt            *time.Time            `json:"completedAt,omitempty"`
}
