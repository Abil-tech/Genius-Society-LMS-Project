package teacherdto

// WalasStudentDetailDTO represents detailed information of a student in the homeroom class.
type WalasStudentDetailDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	NIS          string   `json:"nis"`
	Gender       string   `json:"gender"`
	Email        string   `json:"email,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	ParentName   string   `json:"parentName,omitempty"`
	ParentPhone  string   `json:"parentPhone,omitempty"`
	AverageScore *float64 `json:"averageScore,omitempty"`
	Rank         int      `json:"rank,omitempty"`
}

// WalasSubjectScoreDTO represents score components for one subject.
type WalasSubjectScoreDTO struct {
	SubjectID   string  `json:"subjectId"`
	SubjectName string  `json:"subjectName"`
	DailyAvg    float64 `json:"dailyAvg"`
	UTS         float64 `json:"uts"`
	UAS         float64 `json:"uas"`
	FinalScore  float64 `json:"finalScore"`
}

// WalasGradeRekapRowDTO represents a student row in the Walas full-grade recap table.
type WalasGradeRekapRowDTO struct {
	StudentID      string                 `json:"studentId"`
	StudentName    string                 `json:"studentName"`
	StudentNIS     string                 `json:"studentNis"`
	SubjectScores  []WalasSubjectScoreDTO `json:"subjectScores"`
	OverallAverage float64                `json:"overallAverage"`
	Rank           int                    `json:"rank,omitempty"`
}
