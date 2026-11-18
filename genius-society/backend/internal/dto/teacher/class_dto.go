package teacherdto

// SubjectBrief represents a brief subject representation for class and teacher views.
type SubjectBrief struct {
	ID   string `json:"id"`
	Code string `json:"code,omitempty"`
	Name string `json:"name"`
}

// TeacherClassDTO represents a class accessible to the teacher (teaching or walas).
type TeacherClassDTO struct {
	ID             string         `json:"id"`
	Code           string         `json:"code"`
	Name           string         `json:"name"`
	EducationLevel string         `json:"educationLevel"`
	Department     string         `json:"department,omitempty"`
	AcademicYearID string         `json:"academicYearId,omitempty"`
	StudentCount   int64          `json:"studentCount"`
	Subjects       []SubjectBrief `json:"subjects"`
	IsWalas        bool           `json:"isWalas"`
}

// TeacherStudentDTO represents a student in a teacher's class view.
type TeacherStudentDTO struct {
	ID         string `json:"id"`
	StudentID  string `json:"studentId"`
	Name       string `json:"name"`
	NIS        string `json:"nis"`
	NISN       string `json:"nisn,omitempty"`
	Gender     string `json:"gender,omitempty"`
	ClassID    string `json:"classId"`
	ClassName  string `json:"className"`
	Email      string `json:"email,omitempty"`
	Status     string `json:"status"`
	HasAccount bool   `json:"hasAccount"`
}

// TeacherSubjectDTO represents a subject taught by the teacher.
type TeacherSubjectDTO struct {
	ID       string   `json:"id"`
	Code     string   `json:"code"`
	Name     string   `json:"name"`
	ClassIDs []string `json:"classIds"`
}

// ClassStudentsMeta holds pagination metadata.
type ClassStudentsMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}
