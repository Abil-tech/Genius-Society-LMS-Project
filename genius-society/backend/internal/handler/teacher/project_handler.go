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

// ProjectHandler handles HTTP requests for teacher project endpoints.
type ProjectHandler struct {
	svc teachersvc.ProjectService
}

// NewProjectHandler constructs a new ProjectHandler instance.
func NewProjectHandler(svc teachersvc.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

// ListProjects handles GET /api/teacher/projects.
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	classID := strings.TrimSpace(c.Query("classId"))
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		search = strings.TrimSpace(c.Query("search"))
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	list, meta, err := h.svc.ListProjects(c.Request.Context(), scope, classID, subjectID, search, page, limit)
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

// GetProjectDetail handles GET /api/teacher/projects/:id.
func (h *ProjectHandler) GetProjectDetail(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	prj, err := h.svc.GetProjectDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prj,
	})
}

// CreateProject handles POST /api/teacher/projects.
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	var payload teacherdto.ProjectCreatePayload
	var files []*multipart.FileHeader

	contentType := c.ContentType()
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			handler.RespondError(c, handler.NewValidationError("Gagal membaca multipart form"))
			return
		}

		// Try parsing from a JSON "data" field first (API client approach).
		jsonData := c.PostForm("data")
		if jsonData != "" {
			if err := json.Unmarshal([]byte(jsonData), &payload); err != nil {
				handler.RespondError(c, handler.NewValidationError("Payload JSON pada field data tidak valid"))
				return
			}
		} else {
			// Frontend sends individual form fields — parse them one by one.
			payload.Title = strings.TrimSpace(c.PostForm("title"))
			payload.Description = strings.TrimSpace(c.PostForm("description"))
			payload.Instructions = strings.TrimSpace(c.PostForm("instructions"))
			payload.SubjectID = strings.TrimSpace(c.PostForm("subjectId"))
			if ar := strings.TrimSpace(c.PostForm("allowResubmit")); ar == "true" {
				payload.AllowResubmit = true
			}
			// classIds may appear multiple times (one per append call).
			if form != nil && form.Value != nil {
				payload.ClassIDs = form.Value["classIds"]
			}
			// Parse deadline from ISO 8601 string.
			if dlStr := strings.TrimSpace(c.PostForm("deadline")); dlStr != "" {
				t, err := time.Parse(time.RFC3339, dlStr)
				if err != nil {
					// Try without nanoseconds (e.g. "2026-12-01T00:00:00Z").
					t, err = time.Parse("2006-01-02T15:04:05Z07:00", dlStr)
				}
				if err != nil {
					// Try date-time-local format "2006-01-02T15:04".
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

	res, err := h.svc.CreateProject(c.Request.Context(), scope, &payload, files)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    res,
	})
}

// UpdateProject handles PUT /api/teacher/projects/:id.
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	var payload teacherdto.ProjectUpdatePayload
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
			if v := strings.TrimSpace(c.PostForm("instructions")); v != "" {
				payload.Instructions = v
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

	res, err := h.svc.UpdateProject(c.Request.Context(), scope, id, &payload, files)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// DeleteProject handles DELETE /api/teacher/projects/:id.
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	if err := h.svc.DeleteProject(c.Request.Context(), scope, id); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Projek berhasil dihapus",
	})
}

// DownloadFile handles GET /api/teacher/projects/:id/files/:index/download.
func (h *ProjectHandler) DownloadFile(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
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

// ListSubmissions handles GET /api/teacher/projects/:id/submissions.
func (h *ProjectHandler) ListSubmissions(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
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

// GradeSubmission handles PUT /api/teacher/projects/:id/submissions/:submissionId/grade.
func (h *ProjectHandler) GradeSubmission(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	subIDStr := strings.TrimSpace(c.Param("submissionId"))
	subID, err := bson.ObjectIDFromHex(subIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID submission tidak valid"))
		return
	}

	var payload teacherdto.GradeProjectSubmissionPayload
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

// SubmitProject handles POST /api/teacher/projects/:id/submit.
// Rule: POST /submit adalah milik Murid; Guru 403.
func (h *ProjectHandler) SubmitProject(c *gin.Context) {
	handler.RespondError(c, handler.NewForbiddenError("Pengumpulan projek hanya dapat dilakukan oleh murid"))
}
