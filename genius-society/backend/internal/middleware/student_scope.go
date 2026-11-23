package middleware

import (
	"context"
	"net/http"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	ContextKeyStudentScope = "studentScope"
)

// StudentScope holds the authenticated student's academic context.
type StudentScope struct {
	SchoolID       bson.ObjectID
	StudentID      bson.ObjectID
	UserID         bson.ObjectID
	ClassID        bson.ObjectID
	AcademicYearID bson.ObjectID
	Semester       int
}

// StudentScopeRepository defines methods needed by StudentScopeMiddleware.
type StudentScopeRepository interface {
	FindUserByID(ctx context.Context, id bson.ObjectID) (*model.User, error)
	FindStudentByUserID(ctx context.Context, userID bson.ObjectID) (*model.Student, error)
	FindActiveAcademicYear(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error)
	FindActiveEnrollment(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) (*model.StudentEnrollment, error)
}

// StudentScopeMiddleware enforces student permissions and attaches StudentScope to the gin context.
func StudentScopeMiddleware(repo StudentScopeRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := GetUserID(c)
		if userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "UNAUTHORIZED",
					Message: "Sesi tidak ditemukan atau kedaluwarsa",
				},
			})
			return
		}

		uid, err := bson.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "UNAUTHORIZED",
					Message: "ID pengguna tidak valid",
				},
			})
			return
		}

		ctx := c.Request.Context()

		usr, err := repo.FindUserByID(ctx, uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "INTERNAL_ERROR",
					Message: "Terjadi kesalahan internal pada server",
				},
			})
			return
		}
		if usr == nil || !usr.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "ACCOUNT_INACTIVE",
					Message: "Akun tidak aktif atau dinonaktifkan",
				},
			})
			return
		}

		if usr.Role != model.RoleStudent {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "FORBIDDEN",
					Message: "Akses ditolak: endpoint ini khusus untuk murid",
				},
			})
			return
		}

		if usr.SchoolID.IsZero() {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "FORBIDDEN",
					Message: "Pengguna tidak terdaftar pada sekolah manapun",
				},
			})
			return
		}

		student, err := repo.FindStudentByUserID(ctx, uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "INTERNAL_ERROR",
					Message: "Terjadi kesalahan internal pada server",
				},
			})
			return
		}
		if student == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "FORBIDDEN",
					Message: "Data profil murid tidak ditemukan",
				},
			})
			return
		}

		ay, err := repo.FindActiveAcademicYear(ctx, usr.SchoolID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "INTERNAL_ERROR",
					Message: "Terjadi kesalahan internal pada server",
				},
			})
			return
		}
		if ay == nil {
			c.AbortWithStatusJSON(http.StatusConflict, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "ACADEMIC_YEAR_NOT_ACTIVE",
					Message: "Tidak ada tahun ajaran aktif pada sekolah ini",
				},
			})
			return
		}

		// Determine ClassID: from active enrollment or fallback to student.ClassID
		classID := student.ClassID
		enrollment, err := repo.FindActiveEnrollment(ctx, usr.SchoolID, student.ID, ay.ID)
		if err == nil && enrollment != nil && !enrollment.ClassID.IsZero() {
			classID = enrollment.ClassID
		}

		if classID.IsZero() {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "NO_CLASS_ASSIGNED",
					Message: "Murid belum ditempatkan pada kelas aktif manapun",
				},
			})
			return
		}

		scope := &StudentScope{
			SchoolID:       usr.SchoolID,
			StudentID:      student.ID,
			UserID:         usr.ID,
			ClassID:        classID,
			AcademicYearID: ay.ID,
			Semester:       ay.Semester,
		}

		c.Set(ContextKeyStudentScope, scope)
		c.Set("studentID", student.ID.Hex())
		c.Set("schoolID", usr.SchoolID.Hex())

		c.Next()
	}
}

// GetStudentScope retrieves StudentScope from gin Context.
func GetStudentScope(c *gin.Context) (*StudentScope, bool) {
	v, exists := c.Get(ContextKeyStudentScope)
	if !exists {
		return nil, false
	}
	scope, ok := v.(*StudentScope)
	return scope, ok
}

// MustGetStudentScope retrieves StudentScope from gin Context or panics.
func MustGetStudentScope(c *gin.Context) *StudentScope {
	scope, ok := GetStudentScope(c)
	if !ok || scope == nil {
		panic("StudentScope not found in context")
	}
	return scope
}
