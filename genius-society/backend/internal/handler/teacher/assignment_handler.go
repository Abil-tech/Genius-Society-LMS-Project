package teacherhandler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AssignmentHandler handles HTTP requests for teacher assignment endpoints.
type AssignmentHandler struct {
	svc teachersvc.AssignmentService
}

// NewAssignmentHandler constructs a new AssignmentHandler instance.
func NewAssignmentHandler(svc teachersvc.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{svc: svc}
}

// ListAssignments handles GET /api/teacher/assignments.
func (h *AssignmentHandler) ListAssignments(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	classID := strings.TrimSpace(c.Query("classId"))
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		search = strings.TrimSpace(c.Query("search"))
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	list, meta, err := h.svc.ListAssignments(c.Request.Context(), scope, classID, subjectID, search, page, limit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
		"meta":    meta,
	})
}

// GetAssignmentDetail handles GET /api/teacher/assignments/:id.
func (h *AssignmentHandler) GetAssignmentDetail(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	asg, err := h.svc.GetAssignmentDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    asg,
	})
}

// CreateAssignment handles POST /api/teacher/assignments.
func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	var payload teacherdto.AssignmentCreatePayload
	var files []*multipart.FileHeader

	contentType := c.ContentType()
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			handler.RespondError(c, handler.NewValidationError("Gagal membaca multipart form"))
			return
		}

		jsonData := c.PostForm("data")
		if jsonData != "" {
			if err := json.Unmarshal([]byte(jsonData), &payload); err != nil {
				handler.RespondError(c, handler.NewValidationError("Payload JSON pada field data tidak valid"))
				return
			}
		} else {
			// Frontend sends individual form fields.
			payload.Title = strings.TrimSpace(c.PostForm("title"))
			payload.Description = strings.TrimSpace(c.PostForm("description"))
			payload.SubjectID = strings.TrimSpace(c.PostForm("subjectId"))
			if ar := strings.TrimSpace(c.PostForm("allowResubmit")); ar == "true" {
				payload.AllowResubmit = true
			}
			if form != nil && form.Value != nil {
				payload.ClassIDs = form.Value["classIds"]
			}
			if dlStr := strings.TrimSpace(c.PostForm("deadline")); dlStr != "" {
				t, err := time.Parse(time.RFC3339, dlStr)
				if err != nil {
					t, err = time.Parse("2006-01-02T15:04:05Z07:00", dlStr)
				}
				if err != nil {
					t, err = time.Parse("2006-01-02T15:04", dlStr)
				}
				if err == nil {
					payload.Deadline = t.UTC()
				}
			}
		}

		if form != nil && form.File != nil {
			if fList, ok := form.File["attachments"]; ok {
				files = fList
			} else if fList, ok := form.File["files"]; ok {
				files = fList
			} else if fList, ok := form.File["file"]; ok {
				files = fList
			}
		}
	} else {
		if err := c.ShouldBindJSON(&payload); err != nil {
			handler.BindError(c, err)
			return
		}
	}

	res, err := h.svc.CreateAssignment(c.Request.Context(), scope, &payload, files)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    res,
	})
}

// UpdateAssignment handles PUT /api/teacher/assignments/:id.
func (h *AssignmentHandler) UpdateAssignment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	var payload teacherdto.AssignmentUpdatePayload
	var files []*multipart.FileHeader

	contentType := c.ContentType()
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			handler.RespondError(c, handler.NewValidationError("Gagal membaca multipart form"))
			return
		}

		jsonData := c.PostForm("data")
		if jsonData != "" {
			if err := json.Unmarshal([]byte(jsonData), &payload); err != nil {
				handler.RespondError(c, handler.NewValidationError("Payload JSON pada field data tidak valid"))
				return
			}
		} else {
			// Frontend sends individual form fields.
			if v := strings.TrimSpace(c.PostForm("title")); v != "" {
				payload.Title = v
			}
			if v := strings.TrimSpace(c.PostForm("description")); v != "" {
				payload.Description = v
			}
			if form != nil && form.Value != nil {
				if cids, ok := form.Value["classIds"]; ok && len(cids) > 0 {
					payload.ClassIDs = cids
				}
			}
			if ar := strings.TrimSpace(c.PostForm("allowResubmit")); ar != "" {
				b := ar == "true"
				payload.AllowResubmit = &b
			}
			if dlStr := strings.TrimSpace(c.PostForm("deadline")); dlStr != "" {
				t, err := time.Parse(time.RFC3339, dlStr)
				if err != nil {
					t, err = time.Parse("2006-01-02T15:04:05Z07:00", dlStr)
				}
				if err != nil {
					t, err = time.Parse("2006-01-02T15:04", dlStr)
				}
				if err == nil {
					utc := t.UTC()
					payload.Deadline = &utc
				}
			}
		}

		if form != nil && form.File != nil {
			if fList, ok := form.File["attachments"]; ok {
				files = fList
			} else if fList, ok := form.File["files"]; ok {
				files = fList
			} else if fList, ok := form.File["file"]; ok {
				files = fList
			}
		}
	} else {
		if err := c.ShouldBindJSON(&payload); err != nil {
			handler.BindError(c, err)
			return
		}
	}

	res, err := h.svc.UpdateAssignment(c.Request.Context(), scope, id, &payload, files)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// DeleteAssignment handles DELETE /api/teacher/assignments/:id.
func (h *AssignmentHandler) DeleteAssignment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	if err := h.svc.DeleteAssignment(c.Request.Context(), scope, id); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tugas berhasil dihapus",
	})
}

// DownloadFile handles GET /api/teacher/assignments/:id/files/:index/download.
func (h *AssignmentHandler) DownloadFile(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	idxStr := strings.TrimSpace(c.Param("index"))
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 0 {
		handler.RespondError(c, handler.NewValidationError("Index berkas tidak valid"))
		return
	}

	downloadURL, err := h.svc.GetDownloadURL(c.Request.Context(), scope, id, idx)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url":       downloadURL,
			"expiresIn": 300,
		},
	})
}

// ListSubmissions handles GET /api/teacher/assignments/:id/submissions.
func (h *AssignmentHandler) ListSubmissions(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	status := strings.TrimSpace(c.Query("status"))
	classID := strings.TrimSpace(c.Query("classId"))
	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		search = strings.TrimSpace(c.Query("search"))
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	submissions, meta, err := h.svc.ListSubmissions(c.Request.Context(), scope, id, status, classID, search, page, limit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    submissions,
		"meta":    meta,
	})
}

// GradeSubmission handles PUT /api/teacher/assignments/:id/submissions/:submissionId/grade.
func (h *AssignmentHandler) GradeSubmission(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	subIDStr := strings.TrimSpace(c.Param("submissionId"))
	subID, err := bson.ObjectIDFromHex(subIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID submission tidak valid"))
		return
	}

	var payload teacherdto.GradeSubmissionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.GradeSubmission(c.Request.Context(), scope, id, subID, payload.Score, payload.Feedback); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nilai dan umpan balik berhasil disimpan",
	})
}

// SubmitAssignment handles POST /api/teacher/assignments/:id/submit.
// Rule: POST /submit adalah milik Murid; Guru 403.
func (h *AssignmentHandler) SubmitAssignment(c *gin.Context) {
	handler.RespondError(c, handler.NewForbiddenError("Pengumpulan tugas hanya dapat dilakukan oleh murid"))
}
