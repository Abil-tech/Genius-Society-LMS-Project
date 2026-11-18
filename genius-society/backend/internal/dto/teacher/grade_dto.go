package teacherdto

// StudentGradeRowDTO represents a single student's row in the teacher grades recap table.
type StudentGradeRowDTO struct {
	StudentID       string   `json:"studentId"`
	StudentName     string   `json:"studentName"`
	StudentNIS      string   `json:"studentNis"`
	AssignmentScore *float64 `json:"assignmentScore,omitempty"`
	AssessmentScore *float64 `json:"assessmentScore,omitempty"`
	ProjectScore    *float64 `json:"projectScore,omitempty"`
	DailyAvg        *float64 `json:"dailyAvg,omitempty"`
	UTSScore        *float64 `json:"utsScore,omitempty"`
	UASScore        *float64 `json:"uasScore,omitempty"`
	FinalScore      *float64 `json:"finalScore,omitempty"`
}

// UpdateStudentGradePayload represents the request body to update a student's UTS/UAS score.
type UpdateStudentGradePayload struct {
	StudentID string  `json:"studentId" binding:"required"`
	SubjectID string  `json:"subjectId" binding:"required"`
	ClassID   string  `json:"classId" binding:"required"`
	Type      string  `json:"type" binding:"required,oneof=uts uas"`
	Score     float64 `json:"score" binding:"min=0,max=100"`
}

// GradeItemDTO represents a single grade record response.
type GradeItemDTO struct {
	ID             string  `json:"id"`
	StudentID      string  `json:"studentId"`
	StudentName    string  `json:"studentName,omitempty"`
	StudentNIS     string  `json:"studentNis,omitempty"`
	SubjectID      string  `json:"subjectId"`
	SubjectName    string  `json:"subjectName,omitempty"`
	ClassID        string  `json:"classId"`
	ClassName      string  `json:"className,omitempty"`
	AcademicYearID string  `json:"academicYearId,omitempty"`
	Semester       int     `json:"semester,omitempty"`
	Type           string  `json:"type"`
	SourceTitle    string  `json:"sourceTitle,omitempty"`
	Score          float64 `json:"score"`
	CreatedAt      string  `json:"createdAt"`
}
