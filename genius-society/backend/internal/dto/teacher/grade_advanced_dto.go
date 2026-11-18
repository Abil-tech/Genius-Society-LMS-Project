package teacherdto

// ExamEntryDTO represents an individual student's score in an exam bulk input.
type ExamEntryDTO struct {
	StudentID string  `json:"studentId"`
	Score     float64 `json:"score"`
}

// BulkExamsPayload represents payload for PUT /api/teacher/grades/exams.
type BulkExamsPayload struct {
	ClassID   string         `json:"classId"`
	SubjectID string         `json:"subjectId"`
	Semester  int            `json:"semester"`
	Type      string         `json:"type"` // "sts"
	Entries   []ExamEntryDTO `json:"entries"`
}

// ManualGradePayload represents payload for PUT /api/teacher/{items}/:id/students/:studentId/grade.
type ManualGradePayload struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback,omitempty"`
	Reason   string  `json:"reason,omitempty"`
}

// OutcomeRecapDTO represents learning outcome status for a single student in grade recap.
type OutcomeRecapDTO struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Score            *float64 `json:"score"`
	ProvisionalScore *float64 `json:"provisionalScore"`
	GradedCount      int      `json:"gradedCount"`
	TotalCount       int      `json:"totalCount"`
	ManualCount      int      `json:"manualCount"`
}

// StudentGradeRecapDTO represents an individual student's row in grade recap.
type StudentGradeRecapDTO struct {
	StudentID         string            `json:"studentId"`
	StudentName       string            `json:"studentName"`
	NIS               string            `json:"nis"`
	Outcomes          []OutcomeRecapDTO `json:"outcomes"`
	AvgOutcomes       *float64          `json:"avgOutcomes"`
	STS               *float64          `json:"sts"`
	SemesterStsScore  *float64          `json:"semesterStsScore"`
	IsComplete        bool              `json:"isComplete"`
	MissingComponents []string          `json:"missingComponents"`
}

// GradeRecapResponseDTO represents full class recap response.
type GradeRecapResponseDTO struct {
	ClassID           string                 `json:"classId"`
	ClassName         string                 `json:"className"`
	SubjectID         string                 `json:"subjectId"`
	SubjectName       string                 `json:"subjectName"`
	Semester          int                    `json:"semester"`
	OutcomeWeight     float64                `json:"outcomeWeight"`
	StsWeight         float64                `json:"stsWeight"`
	UnmappedItemCount int                    `json:"unmappedItemCount"`
	EmptyOutcomes     []string               `json:"emptyOutcomes"`
	Students          []StudentGradeRecapDTO `json:"students"`
}
