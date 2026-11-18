package teacherhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geniussociety/backend/internal/dto"
	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockClassService struct {
	getClassesFunc     func(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherClassDTO, error)
	getClassDetailFunc func(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID) (*teacherdto.TeacherClassDTO, error)
	getStudentsFunc    func(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID, search string, page, limit int64) ([]teacherdto.TeacherStudentDTO, *teacherdto.ClassStudentsMeta, error)
	getSubjectsFunc    func(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherSubjectDTO, error)
}

func (m *mockClassService) GetClasses(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherClassDTO, error) {
	if m.getClassesFunc != nil {
		return m.getClassesFunc(ctx, scope)
	}
	return []teacherdto.TeacherClassDTO{}, nil
}

func (m *mockClassService) GetClassDetail(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID) (*teacherdto.TeacherClassDTO, error) {
	if m.getClassDetailFunc != nil {
		return m.getClassDetailFunc(ctx, scope, classID)
	}
	return nil, nil
}

func (m *mockClassService) GetClassStudents(ctx context.Context, scope *middleware.TeacherScope, classID bson.ObjectID, search string, page, limit int64) ([]teacherdto.TeacherStudentDTO, *teacherdto.ClassStudentsMeta, error) {
	if m.getStudentsFunc != nil {
		return m.getStudentsFunc(ctx, scope, classID, search, page, limit)
	}
	return []teacherdto.TeacherStudentDTO{}, &teacherdto.ClassStudentsMeta{}, nil
}

func (m *mockClassService) GetTeacherSubjects(ctx context.Context, scope *middleware.TeacherScope) ([]teacherdto.TeacherSubjectDTO, error) {
	if m.getSubjectsFunc != nil {
		return m.getSubjectsFunc(ctx, scope)
	}
	return []teacherdto.TeacherSubjectDTO{}, nil
}

func setupTestRouter(svc teachersvc.ClassService, scope *middleware.TeacherScope) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyTeacherScope, scope)
		c.Next()
	})

	h := NewClassHandler(svc)
	g := r.Group("/api/teacher")
	{
		g.GET("/classes", h.ListClasses)
		g.GET("/classes/:id", h.GetClassDetail)
		g.GET("/classes/:id/students", h.ListClassStudents)
		g.GET("/subjects", h.ListTeacherSubjects)
	}
	return r
}

func TestHandler_ListClasses(t *testing.T) {
	scope := &middleware.TeacherScope{SchoolID: bson.NewObjectID()}
	mockSvc := &mockClassService{
		getClassesFunc: func(ctx context.Context, sc *middleware.TeacherScope) ([]teacherdto.TeacherClassDTO, error) {
			return []teacherdto.TeacherClassDTO{
				{ID: "class1", Name: "X-IPA-1", StudentCount: 30},
			}, nil
		},
	}

	r := setupTestRouter(mockSvc, scope)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/teacher/classes", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Success bool                         `json:"success"`
		Data    []teacherdto.TeacherClassDTO `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success || len(resp.Data) != 1 || resp.Data[0].Name != "X-IPA-1" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestHandler_GetClassDetail_InvalidID(t *testing.T) {
	scope := &middleware.TeacherScope{SchoolID: bson.NewObjectID()}
	mockSvc := &mockClassService{}
	r := setupTestRouter(mockSvc, scope)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/teacher/classes/invalid-hex", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid ID, got %d", w.Code)
	}

	var errResp dto.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp.Error.Code != "VALIDATION_ERROR" {
		t.Errorf("expected VALIDATION_ERROR, got %s", errResp.Error.Code)
	}
}

func TestHandler_GetClassDetail_NotFound(t *testing.T) {
	scope := &middleware.TeacherScope{SchoolID: bson.NewObjectID()}
	mockSvc := &mockClassService{
		getClassDetailFunc: func(ctx context.Context, sc *middleware.TeacherScope, cID bson.ObjectID) (*teacherdto.TeacherClassDTO, error) {
			return nil, handler.ErrNotFound
		},
	}
	r := setupTestRouter(mockSvc, scope)

	validHex := bson.NewObjectID().Hex()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/teacher/classes/"+validHex, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for forbidden/missing class, got %d", w.Code)
	}
}

func TestHandler_ListClassStudents_Success(t *testing.T) {
	scope := &middleware.TeacherScope{SchoolID: bson.NewObjectID()}
	classID := bson.NewObjectID()

	mockSvc := &mockClassService{
		getStudentsFunc: func(ctx context.Context, sc *middleware.TeacherScope, cID bson.ObjectID, search string, page, limit int64) ([]teacherdto.TeacherStudentDTO, *teacherdto.ClassStudentsMeta, error) {
			return []teacherdto.TeacherStudentDTO{
				{ID: "st1", Name: "Siswa 1", NIS: "1001", HasAccount: true},
			}, &teacherdto.ClassStudentsMeta{Page: 1, Limit: 20, Total: 1, TotalPages: 1}, nil
		},
	}
	r := setupTestRouter(mockSvc, scope)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/teacher/classes/"+classID.Hex()+"/students?page=1&limit=20", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Success bool                           `json:"success"`
		Data    []teacherdto.TeacherStudentDTO `json:"data"`
		Meta    teacherdto.ClassStudentsMeta   `json:"meta"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success || len(resp.Data) != 1 || resp.Meta.Total != 1 {
		t.Errorf("unexpected response: %+v", resp)
	}
}
