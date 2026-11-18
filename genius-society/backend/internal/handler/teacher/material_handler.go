package teacherhandler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MaterialHandler handles HTTP requests for teacher material endpoints.
type MaterialHandler struct {
	svc teachersvc.MaterialService
}

// NewMaterialHandler constructs a new MaterialHandler instance.
func NewMaterialHandler(svc teachersvc.MaterialService) *MaterialHandler {
	return &MaterialHandler{svc: svc}
}

// ListMaterials handles GET /api/teacher/materials.
func (h *MaterialHandler) ListMaterials(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	classID := strings.TrimSpace(c.Query("classId"))
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		search = strings.TrimSpace(c.Query("search"))
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	list, meta, err := h.svc.ListMaterials(c.Request.Context(), scope, classID, subjectID, search, page, limit)
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

// GetMaterialDetail handles GET /api/teacher/materials/:id.
func (h *MaterialHandler) GetMaterialDetail(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID materi tidak valid"))
		return
	}

	mat, err := h.svc.GetMaterialDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mat,
	})
}

// CreateMaterial handles POST /api/teacher/materials.
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	var payload teacherdto.MaterialCreatePayload
	var files []*multipart.FileHeader

	// Support both multipart form (field "data" + files) and application/json
	contentType := c.ContentType()
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			handler.RespondError(c, handler.NewValidationError("Gagal membaca multipart form"))
			return
		}

		jsonData := c.PostForm("data")
		if jsonData == "" {
			// Try fallback direct form values
			payload.Title = c.PostForm("title")
			payload.Description = c.PostForm("description")
			payload.SubjectID = c.PostForm("subjectId")
			payload.Link = c.PostForm("link")
			payload.PublishNow = c.PostForm("publishNow") == "true"
			if classIDsStr := c.PostForm("classIds"); classIDsStr != "" {
				payload.ClassIDs = strings.Split(classIDsStr, ",")
			}
		} else {
			if err := json.Unmarshal([]byte(jsonData), &payload); err != nil {
				handler.RespondError(c, handler.NewValidationError("Payload JSON pada field data tidak valid"))
				return
			}
		}

		if form != nil && form.File != nil {
			if fList, ok := form.File["files"]; ok {
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

	res, err := h.svc.CreateMaterial(c.Request.Context(), scope, &payload, files)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    res,
	})
}

// UpdateMaterial handles PUT /api/teacher/materials/:id.
func (h *MaterialHandler) UpdateMaterial(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID materi tidak valid"))
		return
	}

	var payload teacherdto.MaterialUpdatePayload
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
		}

		if form != nil && form.File != nil {
			if fList, ok := form.File["files"]; ok {
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

	res, err := h.svc.UpdateMaterial(c.Request.Context(), scope, id, &payload, files)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// DeleteMaterial handles DELETE /api/teacher/materials/:id.
func (h *MaterialHandler) DeleteMaterial(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID materi tidak valid"))
		return
	}

	if err := h.svc.DeleteMaterial(c.Request.Context(), scope, id); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Materi berhasil dihapus",
	})
}

// DownloadFile handles GET /api/teacher/materials/:id/files/:index/download.
func (h *MaterialHandler) DownloadFile(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID materi tidak valid"))
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
