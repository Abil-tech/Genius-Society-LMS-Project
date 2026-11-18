package teacherdto

// TeacherScheduleItemDTO represents a teaching schedule slot for the teacher portal.
type TeacherScheduleItemDTO struct {
	ID        string `json:"id"`
	Day       int    `json:"day"` // 1=Monday, ..., 7=Sunday
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Subject   string `json:"subject"`
	ClassName string `json:"className"`
	Room      string `json:"room,omitempty"`
	ClassID   string `json:"classId,omitempty"`
	SubjectID string `json:"subjectId,omitempty"`
}
