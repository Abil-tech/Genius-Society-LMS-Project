package teacherdto

// DashboardKPI holds the aggregate KPI counters for the teacher dashboard.
type DashboardKPI struct {
	ClassesCount           int `json:"classesCount"`
	StudentsCount          int `json:"studentsCount"`
	PendingGradesCount     int `json:"pendingGradesCount"`
	ActiveAssignmentsCount int `json:"activeAssignmentsCount"`
	ActiveAssessmentsCount int `json:"activeAssessmentsCount"`
	ActiveProjectsCount    int `json:"activeProjectsCount"`
}

// DashboardScheduleItem represents a single schedule slot for today.
type DashboardScheduleItem struct {
	ID        string `json:"id"`
	Day       int    `json:"day"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Subject   string `json:"subject"`
	ClassName string `json:"className"`
	Room      string `json:"room,omitempty"`
	ClassID   string `json:"classId,omitempty"`
	SubjectID string `json:"subjectId,omitempty"`
}

// DashboardSubmissionItem represents a recent student submission.
type DashboardSubmissionItem struct {
	ID           string   `json:"id"`
	StudentID    string   `json:"studentId"`
	StudentName  string   `json:"studentName"`
	StudentNIS   string   `json:"studentNis,omitempty"`
	Type         string   `json:"type"` // "assignment" | "project"
	Title        string   `json:"title"`
	SubmittedAt  string   `json:"submittedAt"`
	Graded       bool     `json:"graded"`
	Score        *float64 `json:"score,omitempty"`
	AssignmentID string   `json:"assignmentId,omitempty"`
	ProjectID    string   `json:"projectId,omitempty"`
}

// DashboardWalasClass holds summary info about the walas homeroom class.
type DashboardWalasClass struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Level        string `json:"level"`
	StudentCount int    `json:"studentCount"`
}

// DashboardResponse is the top-level response body for GET /teacher/dashboard.
type DashboardResponse struct {
	KPI               DashboardKPI              `json:"kpi"`
	TodaySchedule     []DashboardScheduleItem   `json:"todaySchedule"`
	RecentSubmissions []DashboardSubmissionItem  `json:"recentSubmissions"`
	IsWalas           bool                      `json:"isWalas"`
	WalasClass        *DashboardWalasClass      `json:"walasClass,omitempty"`
}
