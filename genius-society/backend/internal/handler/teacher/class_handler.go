package teacherhandler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ClassHandler handles HTTP requests for teacher class endpoints.
type ClassHandler struct {
	svc teachersvc.ClassService
}

// NewClassHandler constructs a new ClassHandler instance.
func NewClassHandler(svc teachersvc.ClassService) *ClassHandler {
	return &ClassHandler{svc: svc}
}

// ListClasses handles GET /api/teacher/classes.
func (h *ClassHandler) ListClasses(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	classes, err := h.svc.GetClasses(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    classes,
	})
}

// GetClassDetail handles GET /api/teacher/classes/:id.
func (h *ClassHandler) GetClassDetail(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	idStr := strings.TrimSpace(c.Param("id"))

	classID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID kelas tidak valid"))
		return
	}

	classDetail, err := h.svc.GetClassDetail(c.Request.Context(), scope, classID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    classDetail,
	})
}

// ListClassStudents handles GET /api/teacher/classes/:id/students.
func (h *ClassHandler) ListClassStudents(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	idStr := strings.TrimSpace(c.Param("id"))

	classID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID kelas tidak valid"))
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		search = strings.TrimSpace(c.Query("search"))
	}

	students, meta, err := h.svc.GetClassStudents(c.Request.Context(), scope, classID, search, page, limit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    students,
		"meta":    meta,
	})
}

// ListTeacherSubjects handles GET /api/teacher/subjects.
func (h *ClassHandler) ListTeacherSubjects(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	subjects, err := h.svc.GetTeacherSubjects(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subjects,
	})
}
