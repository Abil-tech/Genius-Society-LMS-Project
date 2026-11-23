package studentdto

import (
	"time"

	"github.com/geniussociety/backend/internal/model"
)

// DeadlineItem represents an upcoming deadline for assignments or projects.
type DeadlineItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // "assignment" | "project"
	Title     string    `json:"title"`
	Subject   string    `json:"subject"`
	Deadline  time.Time `json:"deadline"`
	Submitted bool      `json:"submitted"`
}

// RecentGradeItem represents recently graded tasks for the student.
type RecentGradeItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Subject   string    `json:"subject"`
	Score     float64   `json:"score"`
	GradedAt  time.Time `json:"gradedAt"`
}

// DashboardResponseDTO represents the student dashboard summary.
type DashboardResponseDTO struct {
	ActiveAssignmentsCount int               `json:"activeAssignmentsCount"`
	ActiveAssessmentsCount int               `json:"activeAssessmentsCount"`
	UpcomingDeadlines      []DeadlineItem    `json:"upcomingDeadlines"`
	RecentGrades           []RecentGradeItem `json:"recentGrades"`
	TodaySchedule          []ScheduleItemDTO `json:"todaySchedule"`
}

// MaterialListDTO represents an item in the student materials listing.
type MaterialListDTO struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	SubjectID   string          `json:"subjectId"`
	SubjectName string          `json:"subjectName"`
	TeacherName string          `json:"teacherName"`
	FileCount   int             `json:"fileCount"`
	Link        string          `json:"link,omitempty"`
	PublishAt   *time.Time      `json:"publishAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// MaterialDetailDTO represents the detail of a material with attachments.
type MaterialDetailDTO struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	SubjectID   string          `json:"subjectId"`
	SubjectName string          `json:"subjectName"`
	TeacherName string          `json:"teacherName"`
	Attachments []model.FileRef `json:"attachments"`
	Link        string          `json:"link,omitempty"`
	PublishAt   *time.Time      `json:"publishAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// SubmissionInfoDTO represents a student's submission status and result.
type SubmissionInfoDTO struct {
	ID          string          `json:"id"`
	Submitted   bool            `json:"submitted"`
	SubmittedAt *time.Time      `json:"submittedAt,omitempty"`
	Files       []model.FileRef `json:"files"`
	Note        string          `json:"note,omitempty"`
	Version     int             `json:"version"`
	Score       *float64        `json:"score,omitempty"`
	Feedback    string          `json:"feedback,omitempty"`
	GradedAt    *time.Time      `json:"gradedAt,omitempty"`
	IsLate      bool            `json:"isLate"`
}

// AssignmentListDTO represents an item in student assignment listing.
type AssignmentListDTO struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	SubjectID     string             `json:"subjectId"`
	SubjectName   string             `json:"subjectName"`
	TeacherName   string             `json:"teacherName"`
	Deadline      time.Time          `json:"deadline"`
	AllowResubmit bool               `json:"allowResubmit"`
	Status        string             `json:"status"` // "pending", "submitted", "graded", "overdue"
	Submission    *SubmissionInfoDTO `json:"submission,omitempty"`
	CreatedAt     time.Time          `json:"createdAt"`
}

// AssignmentDetailDTO represents assignment detail with files and student submission.
type AssignmentDetailDTO struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	SubjectID     string             `json:"subjectId"`
	SubjectName   string             `json:"subjectName"`
	TeacherName   string             `json:"teacherName"`
	Deadline      time.Time          `json:"deadline"`
	Attachments   []model.FileRef    `json:"attachments"`
	AllowResubmit bool               `json:"allowResubmit"`
	Status        string             `json:"status"`
	Submission    *SubmissionInfoDTO `json:"submission,omitempty"`
	CreatedAt     time.Time          `json:"createdAt"`
}

// ProjectListDTO represents an item in student project listing.
type ProjectListDTO struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	SubjectID     string             `json:"subjectId"`
	SubjectName   string             `json:"subjectName"`
	TeacherName   string             `json:"teacherName"`
	Deadline      time.Time          `json:"deadline"`
	AllowResubmit bool               `json:"allowResubmit"`
	Status        string             `json:"status"` // "pending", "submitted", "graded", "overdue"
	Submission    *SubmissionInfoDTO `json:"submission,omitempty"`
	CreatedAt     time.Time          `json:"createdAt"`
}

// ProjectDetailDTO represents project detail with instructions, files and submission.
type ProjectDetailDTO struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Instructions  string             `json:"instructions"`
	SubjectID     string             `json:"subjectId"`
	SubjectName   string             `json:"subjectName"`
	TeacherName   string             `json:"teacherName"`
	Deadline      time.Time          `json:"deadline"`
	Attachments   []model.FileRef    `json:"attachments"`
	AllowResubmit bool               `json:"allowResubmit"`
	Status        string             `json:"status"`
	Submission    *SubmissionInfoDTO `json:"submission,omitempty"`
	CreatedAt     time.Time          `json:"createdAt"`
}

// AssessmentListDTO represents an assessment visible to the student.
type AssessmentListDTO struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	DurationMinutes int        `json:"durationMinutes"`
	SubjectID       string     `json:"subjectId"`
	SubjectName     string     `json:"subjectName"`
	TeacherName     string     `json:"teacherName"`
	Category        string     `json:"category"`
	StartAt         *time.Time `json:"startAt,omitempty"`
	EndAt           *time.Time `json:"endAt,omitempty"`
	Status          string     `json:"status"` // "available", "in_progress", "submitted", "completed", "closed"
	AttemptID       string     `json:"attemptId,omitempty"`
	Score           *float64   `json:"score,omitempty"`
}

// QuestionItemDTO represents a question sent to the student during exam.
// CorrectOption is intentionally omitted so students cannot inspect client-side answers!
type QuestionItemDTO struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"` // "mcq" | "essay"
	Text    string   `json:"text"`
	Options []string `json:"options,omitempty"`
	Weight  float64  `json:"weight"`
	Order   int      `json:"order"`
	Answer  string   `json:"answer,omitempty"` // existing answer if saved
}

// AssessmentDetailDTO represents assessment detail before or during start.
type AssessmentDetailDTO struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	DurationMinutes int        `json:"durationMinutes"`
	SubjectID       string     `json:"subjectId"`
	SubjectName     string     `json:"subjectName"`
	TeacherName     string     `json:"teacherName"`
	Category        string     `json:"category"`
	StartAt         *time.Time `json:"startAt,omitempty"`
	EndAt           *time.Time `json:"endAt,omitempty"`
	QuestionCount   int        `json:"questionCount"`
	AttemptStatus   string     `json:"attemptStatus"` // "not_started", "in_progress", "submitted", "graded"
	AttemptID       string     `json:"attemptId,omitempty"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty"`
	TotalScore      *float64   `json:"totalScore,omitempty"`
	ShowResult      bool       `json:"showResult"`
}

// AssessmentStartResponseDTO is returned when a student starts an assessment.
type AssessmentStartResponseDTO struct {
	AttemptID       string            `json:"attemptId"`
	AssessmentID    string            `json:"assessmentId"`
	Title           string            `json:"title"`
	DurationMinutes int               `json:"durationMinutes"`
	StartedAt       time.Time         `json:"startedAt"`
	ExpiresAt       time.Time         `json:"expiresAt"`
	Questions       []QuestionItemDTO `json:"questions"`
}

// AssessmentAnswerSubmission represents one answer submitted by student.
type AssessmentAnswerSubmission struct {
	QuestionID string `json:"questionId"`
	Answer     string `json:"answer"`
}

// AssessmentSubmitPayload is the body payload for POST /assessments/:id/submit
type AssessmentSubmitPayload struct {
	Answers []AssessmentAnswerSubmission `json:"answers"`
}

// AnswerResultDTO represents question review result if showResult is permitted.
type AnswerResultDTO struct {
	QuestionID    string   `json:"questionId"`
	Type          string   `json:"type"`
	Text          string   `json:"text"`
	Options       []string `json:"options,omitempty"`
	StudentAnswer string   `json:"studentAnswer"`
	CorrectOption *int     `json:"correctOption,omitempty"`
	IsCorrect     *bool    `json:"isCorrect,omitempty"`
	Score         float64  `json:"score"`
	MaxWeight     float64  `json:"maxWeight"`
	Feedback      string   `json:"feedback,omitempty"`
}

// AssessmentResultResponseDTO represents full results of a completed assessment.
type AssessmentResultResponseDTO struct {
	AssessmentID string            `json:"assessmentId"`
	Title        string            `json:"title"`
	SubjectName  string            `json:"subjectName"`
	TotalScore   float64           `json:"totalScore"`
	AutoScore    float64           `json:"autoScore"`
	ManualScore  float64           `json:"manualScore"`
	Status       string            `json:"status"` // "submitted", "graded"
	SubmittedAt  *time.Time        `json:"submittedAt,omitempty"`
	Answers      []AnswerResultDTO `json:"answers,omitempty"`
}

// SubjectGradeDTO represents the grade breakdown for a single subject.
type SubjectGradeDTO struct {
	SubjectID          string   `json:"subjectId"`
	SubjectName        string   `json:"subjectName"`
	SubjectCode        string   `json:"subjectCode"`
	TeacherName        string   `json:"teacherName"`
	AssignmentAverage  *float64 `json:"assignmentAverage"`
	AssessmentAverage  *float64 `json:"assessmentAverage"`
	ProjectAverage     *float64 `json:"projectAverage"`
	DailyAverage       *float64 `json:"dailyAverage"` // Tugas + Harian (60%)
	UTSScore           *float64 `json:"utsScore"`     // UTS (20%)
	UASScore           *float64 `json:"uasScore"`     // UAS (20%)
	FinalScore         *float64 `json:"finalScore"`
	GradeLetter        string   `json:"gradeLetter"` // A, B, C, D
}

// GradeDetailItem represents an individual graded component.
type GradeDetailItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // "assignment", "assessment", "project", "uts", "uas"
	Title     string    `json:"title"`
	Score     float64   `json:"score"`
	GradedAt  time.Time `json:"gradedAt"`
	Feedback  string    `json:"feedback,omitempty"`
}

// SubjectGradeDetailDTO represents detailed breakdown of all graded items in a subject.
type SubjectGradeDetailDTO struct {
	SubjectID   string            `json:"subjectId"`
	SubjectName string            `json:"subjectName"`
	FinalScore  *float64          `json:"finalScore"`
	GradeLetter string            `json:"gradeLetter"`
	Assignments []GradeDetailItem `json:"assignments"`
	Assessments []GradeDetailItem `json:"assessments"`
	Projects    []GradeDetailItem `json:"projects"`
	UTS         *GradeDetailItem  `json:"uts,omitempty"`
	UAS         *GradeDetailItem  `json:"uas,omitempty"`
}

// ScheduleItemDTO represents a student's class schedule item.
type ScheduleItemDTO struct {
	ID        string `json:"id"`
	Day       int    `json:"day"`
	DayName   string `json:"dayName"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Subject   string `json:"subject"`
	Teacher   string `json:"teacher"`
	Room      string `json:"room"`
}

// RankingItemDTO represents one student in the class ranking leaderboard.
type RankingItemDTO struct {
	Rank       int     `json:"rank"`
	StudentID  string  `json:"studentId"`
	Name       string  `json:"name"`
	NIS        string  `json:"nis"`
	FinalScore float64 `json:"finalScore"`
	IsSelf     bool    `json:"isSelf"`
}

// ClassRankingResponseDTO contains the leaderboard for the student's class.
type ClassRankingResponseDTO struct {
	ClassName string           `json:"className"`
	MyRank    int              `json:"myRank"`
	Total     int              `json:"totalStudents"`
	Rankings  []RankingItemDTO `json:"rankings"`
}
