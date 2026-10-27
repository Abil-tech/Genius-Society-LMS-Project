package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// mockUserRepo provides a stub for FindByID used by SetSchoolID middleware.
type mockUserRepo struct {
	user *model.User
	err  error
}

func (m *mockUserRepo) FindByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	return m.user, m.err
}

// mockClassRepository implements only Create for service tests.
type mockClassRepository struct {
	receivedSchoolID   bson.ObjectID
	receivedAcademicID bson.ObjectID
	receivedReq        *dto.ClassCreateDTO
	err                error
	classes            map[string]*model.Class
}

func newMockClassRepository() *mockClassRepository {
	return &mockClassRepository{
		classes: make(map[string]*model.Class),
	}
}

func (m *mockClassRepository) Create(ctx context.Context, class *model.Class) error {
	m.receivedSchoolID = class.SchoolID
	m.receivedAcademicID = class.AcademicYearID
	m.receivedReq = &dto.ClassCreateDTO{
		Code:           class.Code,
		Name:           class.Name,
		EducationLevel: class.EducationLevel,
		Department:     class.Department,
		Capacity:       class.Capacity,
		Status:         class.Status,
	}
	if m.err != nil {
		return m.err
	}
	if class.ID.IsZero() {
		class.ID = bson.NewObjectID()
	}
	if m.classes == nil {
		m.classes = make(map[string]*model.Class)
	}
	m.classes[class.ID.Hex()] = class
	return nil
}

func (m *mockClassRepository) ExistsByCode(ctx context.Context, schoolID, academicYearID bson.ObjectID, code string) (bool, error) {
	for _, c := range m.classes {
		if c.SchoolID == schoolID && c.AcademicYearID == academicYearID && c.Code == code && c.Status != "archived" {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockClassRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Class, error) {
	if m.classes == nil {
		return nil, nil
	}
	c, ok := m.classes[id.Hex()]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockClassRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	if m.classes == nil {
		return nil
	}
	c, ok := m.classes[id.Hex()]
	if !ok {
		return nil
	}
	if n, ok := update["name"].(string); ok {
		c.Name = n
	}
	return nil
}

func (m *mockClassRepository) Archive(ctx context.Context, id bson.ObjectID) error {
	if m.classes == nil {
		return nil
	}
	c, ok := m.classes[id.Hex()]
	if !ok {
		return nil
	}
	c.Status = "archived"
	return nil
}

func (m *mockClassRepository) List(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Class, int64, error) {
	var res []*model.Class
	schoolID, _ := filter["schoolId"].(bson.ObjectID)
	for _, c := range m.classes {
		if !schoolID.IsZero() && c.SchoolID != schoolID {
			continue
		}
		res = append(res, c)
	}
	return res, int64(len(res)), nil
}

func (m *mockClassRepository) CountByFilter(ctx context.Context, filter bson.M) (int64, error) {
	list, total, _ := m.List(ctx, filter, 1, 100)
	_ = list
	return total, nil
}

// helper constructs a Gin engine with auth simulation, SetSchoolID and the class handler.
func setupRouter(userRepo interface {
	FindByID(context.Context, bson.ObjectID) (*model.User, error)
}, svc *service.ClassService, role model.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Simulated Auth middleware – inject userID and role.
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, "5f1d7f3e2c2e123456789012")
		c.Set(middleware.ContextKeyRole, role)
		c.Next()
	})
	// Use the SetSchoolID middleware with the provided repo (mockUserRepo implements the required interface).
	r.Use(middleware.SetSchoolID(userRepo))
	classHandler := NewClassManagementHandler(svc)
	classes := r.Group("/classes")
	classes.Use(middleware.RequireRoles(model.RoleAdmin))
	{
		classes.GET("", classHandler.List)
		classes.GET("/:id", classHandler.GetOne)
		classes.POST("", classHandler.Create)
		classes.PUT("/:id", classHandler.Update)
		classes.DELETE("/:id", classHandler.Archive)
	}
	return r
}

func TestSetSchoolID_Success(t *testing.T) {
	user := &model.User{ID: bson.NewObjectID(), SchoolID: bson.NewObjectID()}
	repo := &mockUserRepo{user: user}
	// Simple handler to verify schoolID is set.
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, user.ID.Hex())
		c.Next()
	})
	r.Use(middleware.SetSchoolID(repo))
	r.GET("/test", func(c *gin.Context) {
		if v, ok := c.Get("schoolID"); !ok || v != user.SchoolID.Hex() {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSetSchoolID_NoSchool(t *testing.T) {
	user := &model.User{ID: bson.NewObjectID()}
	repo := &mockUserRepo{user: user}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, user.ID.Hex())
		c.Next()
	})
	r.Use(middleware.SetSchoolID(repo))
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when user has no school, got %d", w.Code)
	}
}

func TestCreateClass_IgnoresHeaderSchoolID(t *testing.T) {
	adminUser := &model.User{ID: bson.NewObjectID(), SchoolID: bson.NewObjectID()}
	repoUser := &mockUserRepo{user: adminUser}
	classRepo := &mockClassRepository{}
	svc := service.NewClassService(classRepo)
	router := setupRouter(repoUser, svc, model.RoleAdmin)

	// Build request body.
	payload := dto.ClassCreateDTO{Code: "CLS01", Name: "Test Class", AcademicYearID: "5f1d7f3e2c2e123456789012", Capacity: 30, Status: "active"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/classes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// Attempt to override school via header – should be ignored.
	req.Header.Set("X-School-ID", "5f1d7f3e2c2e999999999999")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if classRepo.receivedSchoolID != adminUser.SchoolID {
		t.Fatalf("service received wrong schoolID: got %s, want %s", classRepo.receivedSchoolID.Hex(), adminUser.SchoolID.Hex())
	}
}

func TestRequireRoles_AdminOnly(t *testing.T) {
	// Teacher role should be forbidden.
	teacherUser := &model.User{ID: bson.NewObjectID(), SchoolID: bson.NewObjectID()}
	repoUser := &mockUserRepo{user: teacherUser}
	classRepo := &mockClassRepository{}
	svc := service.NewClassService(classRepo)
	router := setupRouter(repoUser, svc, model.RoleTeacher)
	payload := dto.ClassCreateDTO{Code: "CLS01", Name: "Test", AcademicYearID: "5f1d7f3e2c2e123456789012", Capacity: 30, Status: "active"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/classes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-admin role, got %d", w.Code)
	}
}

func TestClassManagementHandler_RBAC_SuperAdminAndPrincipalBlocked(t *testing.T) {
	// SuperAdmin is blocked from academic Class Management CRUD
	superAdmin := &model.User{ID: bson.NewObjectID(), SchoolID: bson.NewObjectID()}
	repoUser := &mockUserRepo{user: superAdmin}
	classRepo := &mockClassRepository{}
	svc := service.NewClassService(classRepo)
	router := setupRouter(repoUser, svc, model.RoleSuperAdmin)

	// Attempt POST /classes
	payload := dto.ClassCreateDTO{Code: "CLS01", Name: "Test", AcademicYearID: "5f1d7f3e2c2e123456789012", Capacity: 30, Status: "active"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/classes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for SuperAdmin on academic Class CRUD, got %d", w.Code)
	}

	// Kepala Sekolah (RolePrincipal) is blocked from POST/PUT/DELETE
	principalUser := &model.User{ID: bson.NewObjectID(), SchoolID: bson.NewObjectID()}
	principalRepo := &mockUserRepo{user: principalUser}
	routerPrincipal := setupRouter(principalRepo, svc, model.RolePrincipal)

	reqP, _ := http.NewRequest(http.MethodPost, "/classes", bytes.NewReader(body))
	reqP.Header.Set("Content-Type", "application/json")
	wP := httptest.NewRecorder()
	routerPrincipal.ServeHTTP(wP, reqP)
	if wP.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for Principal on class creation, got %d", wP.Code)
	}
}

func TestClassManagementHandler_FullCRUD_AdminSuccess(t *testing.T) {
	adminSchoolID := bson.NewObjectID()
	classID := bson.NewObjectID()
	adminUser := &model.User{ID: bson.NewObjectID(), SchoolID: adminSchoolID}
	repoUser := &mockUserRepo{user: adminUser}

	classRepo := &mockClassRepository{}
	// Preload a class in mock repository
	classRepo.Create(context.Background(), &model.Class{
		ID:             classID,
		SchoolID:       adminSchoolID,
		AcademicYearID: bson.NewObjectID(),
		Code:           "7A",
		Name:           "Kelas 7A",
		EducationLevel: "SMP",
		Capacity:       32,
		Status:         "active",
	})
	svc := service.NewClassService(classRepo)
	router := setupRouter(repoUser, svc, model.RoleAdmin)

	// 1. GET /classes
	reqList, _ := http.NewRequest(http.MethodGet, "/classes?page=1&limit=10", nil)
	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 on List, got %d", wList.Code)
	}

	// 2. GET /classes/:id
	reqGet, _ := http.NewRequest(http.MethodGet, "/classes/"+classID.Hex(), nil)
	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 on GetOne, got %d", wGet.Code)
	}

	// 3. PUT /classes/:id
	updatePayload := dto.ClassUpdateDTO{Name: new(string)}
	*updatePayload.Name = "Kelas 7A Updated"
	updBody, _ := json.Marshal(updatePayload)
	reqPut, _ := http.NewRequest(http.MethodPut, "/classes/"+classID.Hex(), bytes.NewReader(updBody))
	reqPut.Header.Set("Content-Type", "application/json")
	wPut := httptest.NewRecorder()
	router.ServeHTTP(wPut, reqPut)
	if wPut.Code != http.StatusOK {
		t.Fatalf("expected 200 on Update, got %d", wPut.Code)
	}

	// 4. DELETE /classes/:id
	reqDel, _ := http.NewRequest(http.MethodDelete, "/classes/"+classID.Hex(), nil)
	wDel := httptest.NewRecorder()
	router.ServeHTTP(wDel, reqDel)
	if wDel.Code != http.StatusOK {
		t.Fatalf("expected 200 on Archive, got %d", wDel.Code)
	}
}
