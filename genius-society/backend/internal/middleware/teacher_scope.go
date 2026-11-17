package middleware

import (
	"net/http"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	ContextKeyTeacherScope    = "teacherScope"
	ContextKeyTeacherID       = "teacherID"
	ContextKeyAcademicYearID  = "academicYearID"
)

// TeacherScope holds the request-scoped teaching authorization context.
type TeacherScope struct {
	SchoolID       bson.ObjectID
	TeacherID      bson.ObjectID
	UserID         bson.ObjectID
	AcademicYearID bson.ObjectID
	Semester       int
	TeachingPairs  []repository.ClassSubjectPair
	WalasClassIDs  []bson.ObjectID
	IsWalas        bool
}

// HasClassSubject checks if the teacher is assigned to teach the specified class and subject.
func (s *TeacherScope) HasClassSubject(classID, subjectID bson.ObjectID) bool {
	for _, p := range s.TeachingPairs {
		if p.ClassID == classID && p.SubjectID == subjectID {
			return true
		}
	}
	return false
}

// HasClass checks if the teacher teaches the specified class for any subject.
func (s *TeacherScope) HasClass(classID bson.ObjectID) bool {
	for _, p := range s.TeachingPairs {
		if p.ClassID == classID {
			return true
		}
	}
	return false
}

// HasSubject checks if the teacher teaches the specified subject in any class.
func (s *TeacherScope) HasSubject(subjectID bson.ObjectID) bool {
	for _, p := range s.TeachingPairs {
		if p.SubjectID == subjectID {
			return true
		}
	}
	return false
}

// IsWalasOf checks if the teacher is walas (homeroom teacher) of the specified class.
func (s *TeacherScope) IsWalasOf(classID bson.ObjectID) bool {
	for _, id := range s.WalasClassIDs {
		if id == classID {
			return true
		}
	}
	return false
}

// CanAccessClass checks if the teacher either teaches this class or is walas of it.
func (s *TeacherScope) CanAccessClass(classID bson.ObjectID) bool {
	return s.HasClass(classID) || s.IsWalasOf(classID)
}

// GetTeachingClassIDs returns unique class IDs taught by this teacher.
func (s *TeacherScope) GetTeachingClassIDs() []bson.ObjectID {
	seen := make(map[bson.ObjectID]bool)
	var list []bson.ObjectID
	for _, p := range s.TeachingPairs {
		if !seen[p.ClassID] {
			seen[p.ClassID] = true
			list = append(list, p.ClassID)
		}
	}
	return list
}

// GetTeachingSubjectIDs returns unique subject IDs taught by this teacher.
func (s *TeacherScope) GetTeachingSubjectIDs() []bson.ObjectID {
	seen := make(map[bson.ObjectID]bool)
	var list []bson.ObjectID
	for _, p := range s.TeachingPairs {
		if !seen[p.SubjectID] {
			seen[p.SubjectID] = true
			list = append(list, p.SubjectID)
		}
	}
	return list
}

// TeacherScopeMiddleware enforces teacher permissions and builds the TeacherScope on every request.
func TeacherScopeMiddleware(repo repository.TeacherScopeRepository) gin.HandlerFunc {
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

		// 1. Check user status from database on every request
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

		if usr.Role != model.RoleTeacher {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "FORBIDDEN",
					Message: "Akses ditolak: endpoint ini khusus untuk guru",
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

		// 2. Load teacher record
		teacher, err := repo.FindTeacherByUserID(ctx, uid)
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
		if teacher == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "FORBIDDEN",
					Message: "Data profil guru tidak ditemukan",
				},
			})
			return
		}

		// 3. Load active academic year
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

		// 4. Load teaching pairs from teacher_classes
		pairs, err := repo.FindTeacherClassPairs(ctx, usr.SchoolID, teacher.ID, ay.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Error: dto.ErrorBody{
					Code:    "INTERNAL_ERROR",
					Message: "Gagal memuat penugasan kelas guru",
				},
			})
			return
		}

		// 5. Determine Walas status per analysis 1c
		isWalas := len(teacher.WalasClassIDs) > 0
		if teacher.IsWalas && len(teacher.WalasClassIDs) > 0 {
			isWalas = true
		}

		scope := &TeacherScope{
			SchoolID:       usr.SchoolID,
			TeacherID:      teacher.ID,
			UserID:         usr.ID,
			AcademicYearID: ay.ID,
			Semester:       ay.Semester,
			TeachingPairs:  pairs,
			WalasClassIDs:  teacher.WalasClassIDs,
			IsWalas:        isWalas,
		}

		c.Set(ContextKeyTeacherScope, scope)
		c.Set(ContextKeyTeacherID, teacher.ID)
		c.Set(ContextKeyAcademicYearID, ay.ID)
		c.Set("schoolID", usr.SchoolID.Hex())

		c.Next()
	}
}

// GetTeacherScope extracts TeacherScope from gin context.
func GetTeacherScope(c *gin.Context) (*TeacherScope, bool) {
	v, exists := c.Get(ContextKeyTeacherScope)
	if !exists {
		return nil, false
	}
	scope, ok := v.(*TeacherScope)
	return scope, ok
}

// MustGetTeacherScope extracts TeacherScope or panics if not present.
func MustGetTeacherScope(c *gin.Context) *TeacherScope {
	scope, ok := GetTeacherScope(c)
	if !ok || scope == nil {
		panic("TeacherScope not found in context")
	}
	return scope
}
