package dto

import "github.com/geniussociety/backend/internal/model"

// Auth
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`    // email or username
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token,omitempty"` // only if not cookie-mode
}

// User
type UserResponse struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Username string     `json:"username"`
	Role     model.Role `json:"role"`
	IsActive bool       `json:"isActive"`
}

type CreateUserRequest struct {
	Name     string     `json:"name" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Username string     `json:"username" binding:"required"`
	Password string     `json:"password" binding:"required,min=8"`
	Role     model.Role `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"omitempty,min=8"`
	IsActive *bool  `json:"isActive"`
}

// Pagination
type PaginationMeta struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
	Total int64 `json:"total"`
}

// API response wrappers
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Admin Dashboard DTO
type AdminDashboardResponse struct {
	AcademicYear struct {
		Name     string `json:"name"`
		Semester int    `json:"semester"`
		IsActive bool   `json:"isActive"`
	} `json:"academicYear"`
	CurrentDate string `json:"currentDate"`
	KPI         struct {
		TotalStudents int64 `json:"totalStudents"`
		TotalTeachers int64 `json:"totalTeachers"`
		TotalClasses  int64 `json:"totalClasses"`
		TotalAccounts int64 `json:"totalAccounts"`
	} `json:"kpi"`
	ClassDistribution []ClassStudentCount `json:"classDistribution"`
	DataCompleteness  struct {
		StudentsWithoutClass int64 `json:"studentsWithoutClass"`
		ClassesWithoutWalas  int64 `json:"classesWithoutWalas"`
		TeachersWithoutMapel int64 `json:"teachersWithoutMapel"`
		InactiveAccounts     int64 `json:"inactiveAccounts"`
	} `json:"dataCompleteness"`
	RecentUsers []RecentUserItem `json:"recentUsers"`
}

type ClassStudentCount struct {
	ClassID   string `json:"classId"`
	ClassName string `json:"className"`
	Count     int64  `json:"count"`
}

type RecentUserItem struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      model.Role `json:"role"`
	CreatedAt string     `json:"createdAt"`
	IsActive  bool       `json:"isActive"`
}

// Student DTOs
type StudentItemResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"userId,omitempty"`
	Name           string `json:"name"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	NIS            string `json:"nis"`
	NISN           string `json:"nisn,omitempty"`
	Gender         string `json:"gender,omitempty"`
	BirthPlace     string `json:"birthPlace,omitempty"`
	BirthDate      string `json:"birthDate,omitempty"`
	Address        string `json:"address,omitempty"`
	ParentName     string `json:"parentName,omitempty"`
	ParentPhone    string `json:"parentPhone,omitempty"`
	EnrollmentYear int    `json:"enrollmentYear,omitempty"`
	PhotoURL       string `json:"photoUrl,omitempty"`
	ClassID        string `json:"classId,omitempty"`
	ClassName      string `json:"className,omitempty"`
	Status         string `json:"status"`
	HasAccount     bool   `json:"hasAccount"`
	IsActive       bool   `json:"isActive"`
	CreatedAt      string `json:"createdAt"`
}

type CreateStudentRequest struct {
	Name           string `json:"name" binding:"required"`
	NIS            string `json:"nis" binding:"required"`
	NISN           string `json:"nisn"`
	Gender         string `json:"gender"`
	BirthPlace     string `json:"birthPlace"`
	BirthDate      string `json:"birthDate"`
	Address        string `json:"address"`
	ParentName     string `json:"parentName"`
	ParentPhone    string `json:"parentPhone"`
	EnrollmentYear int    `json:"enrollmentYear"`
	ClassID        string `json:"classId"`
	Status         string `json:"status"`
	// Opsi buat akun sekaligus
	CreateAccount bool   `json:"createAccount"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
}

type UpdateStudentRequest struct {
	Name           string `json:"name"`
	NIS            string `json:"nis"`
	NISN           string `json:"nisn"`
	Gender         string `json:"gender"`
	BirthPlace     string `json:"birthPlace"`
	BirthDate      string `json:"birthDate"`
	Address        string `json:"address"`
	ParentName     string `json:"parentName"`
	ParentPhone    string `json:"parentPhone"`
	EnrollmentYear int    `json:"enrollmentYear"`
	ClassID        string `json:"classId"`
	Status         string `json:"status"`
}

type AssignClassRequest struct {
	StudentIDs []string `json:"studentIds" binding:"required"`
	ClassID    string   `json:"classId" binding:"required"`
}

type BulkStatusRequest struct {
	StudentIDs []string `json:"studentIds" binding:"required"`
	Status     string   `json:"status" binding:"required"`
}

type CreateStudentAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type StudentStatsResponse struct {
	TotalStudents        int64 `json:"totalStudents"`
	ActiveStudents       int64 `json:"activeStudents"`
	StudentsWithoutClass int64 `json:"studentsWithoutClass"`
	StudentsWithoutAccount int64 `json:"studentsWithoutAccount"`
}

// Teacher DTOs
type TeacherItemResponse struct {
	ID            string   `json:"id"`
	UserID        string   `json:"userId,omitempty"`
	Name          string   `json:"name"`
	Username      string   `json:"username,omitempty"`
	Email         string   `json:"email,omitempty"`
	NIP           string   `json:"nip"`
	Gender        string   `json:"gender,omitempty"`
	Phone         string   `json:"phone,omitempty"`
	Address       string   `json:"address,omitempty"`
	Subject       string   `json:"subject,omitempty"`
	Status        string   `json:"status"`
	IsWalas       bool     `json:"isWalas"`
	WalasClasses  []string `json:"walasClasses,omitempty"`
	HasAccount    bool     `json:"hasAccount"`
	IsActive      bool     `json:"isActive"`
	CreatedAt     string   `json:"createdAt"`
}

type CreateTeacherRequest struct {
	Name          string   `json:"name" binding:"required"`
	NIP           string   `json:"nip" binding:"required"`
	Gender        string   `json:"gender"`
	Phone         string   `json:"phone"`
	Address       string   `json:"address"`
	Subject       string   `json:"subject"`
	Status        string   `json:"status"`
	IsWalas       bool     `json:"isWalas"`
	WalasClassIDs []string `json:"walasClassIds"`
	// Opsi buat akun sekaligus
	CreateAccount bool   `json:"createAccount"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
}

type UpdateTeacherRequest struct {
	Name          string   `json:"name"`
	NIP           string   `json:"nip"`
	Gender        string   `json:"gender"`
	Phone         string   `json:"phone"`
	Address       string   `json:"address"`
	Subject       string   `json:"subject"`
	Status        string   `json:"status"`
	IsWalas       bool     `json:"isWalas"`
	WalasClassIDs []string `json:"walasClassIds"`
}

type CreateTeacherAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type TeacherStatsResponse struct {
	TotalTeachers        int64 `json:"totalTeachers"`
	ActiveTeachers       int64 `json:"activeTeachers"`
	WalasCount           int64 `json:"walasCount"`
	TeachersWithoutMapel int64 `json:"teachersWithoutMapel"`
}

// Subject DTOs
type SubjectItemResponse struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	GradeLevel   string `json:"gradeLevel"`
	KKM          int    `json:"kkm"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	TeacherCount int64  `json:"teacherCount"`
	CreatedAt    string `json:"createdAt"`
}

type CreateSubjectRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category"`
	GradeLevel  string `json:"gradeLevel"`
	KKM         int    `json:"kkm"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateSubjectRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	GradeLevel  string `json:"gradeLevel"`
	KKM         int    `json:"kkm"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type SubjectStatsResponse struct {
	TotalSubjects    int64 `json:"totalSubjects"`
	ActiveSubjects   int64 `json:"activeSubjects"`
	GeneralSubjects  int64 `json:"generalSubjects"`
	SpecialtySubjects int64 `json:"specialtySubjects"`
}
