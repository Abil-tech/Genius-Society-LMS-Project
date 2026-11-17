package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockTeacherScopeRepo struct {
	userFunc          func(ctx context.Context, id bson.ObjectID) (*model.User, error)
	teacherFunc       func(ctx context.Context, id bson.ObjectID) (*model.Teacher, error)
	ayFunc            func(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error)
	pairsFunc         func(ctx context.Context, schoolID, teacherID, ayID bson.ObjectID) ([]repository.ClassSubjectPair, error)
	teacherSubjFunc   func(ctx context.Context, schoolID, teacherID bson.ObjectID) ([]bson.ObjectID, error)
}

func (m *mockTeacherScopeRepo) FindUserByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	if m.userFunc != nil {
		return m.userFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTeacherScopeRepo) FindTeacherByUserID(ctx context.Context, id bson.ObjectID) (*model.Teacher, error) {
	if m.teacherFunc != nil {
		return m.teacherFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTeacherScopeRepo) FindActiveAcademicYear(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error) {
	if m.ayFunc != nil {
		return m.ayFunc(ctx, schoolID)
	}
	return nil, nil
}

func (m *mockTeacherScopeRepo) FindTeacherClassPairs(ctx context.Context, schoolID, teacherID, ayID bson.ObjectID) ([]repository.ClassSubjectPair, error) {
	if m.pairsFunc != nil {
		return m.pairsFunc(ctx, schoolID, teacherID, ayID)
	}
	return nil, nil
}

func (m *mockTeacherScopeRepo) FindTeacherSubjectIDs(ctx context.Context, schoolID, teacherID bson.ObjectID) ([]bson.ObjectID, error) {
	if m.teacherSubjFunc != nil {
		return m.teacherSubjFunc(ctx, schoolID, teacherID)
	}
	return nil, nil
}

func TestTeacherScope_HelperMethods(t *testing.T) {
	class1 := bson.NewObjectID()
	class2 := bson.NewObjectID()
	classWalas := bson.NewObjectID()
	subj1 := bson.NewObjectID()
	subj2 := bson.NewObjectID()

	scope := &TeacherScope{
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: class1, SubjectID: subj1},
			{ClassID: class2, SubjectID: subj2},
		},
		WalasClassIDs: []bson.ObjectID{classWalas},
		IsWalas:       true,
	}

	if !scope.HasClassSubject(class1, subj1) {
		t.Errorf("expected HasClassSubject true for class1, subj1")
	}
	if scope.HasClassSubject(class1, subj2) {
		t.Errorf("expected HasClassSubject false for class1, subj2")
	}
	if !scope.HasClass(class1) || !scope.HasClass(class2) {
		t.Errorf("expected HasClass true for class1 and class2")
	}
	if scope.HasClass(classWalas) {
		t.Errorf("expected HasClass false for classWalas")
	}
	if !scope.IsWalasOf(classWalas) {
		t.Errorf("expected IsWalasOf true for classWalas")
	}
	if !scope.CanAccessClass(classWalas) {
		t.Errorf("expected CanAccessClass true for classWalas")
	}
	if !scope.CanAccessClass(class1) {
		t.Errorf("expected CanAccessClass true for class1")
	}
}

func TestTeacherScopeMiddleware_InactiveUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uid := bson.NewObjectID()
	repo := &mockTeacherScopeRepo{
		userFunc: func(ctx context.Context, id bson.ObjectID) (*model.User, error) {
			return &model.User{
				ID:       uid,
				Role:     model.RoleTeacher,
				IsActive: false, // inactive!
			}, nil
		},
	}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, uid.Hex())
		c.Next()
	})
	r.Use(TeacherScopeMiddleware(repo))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for inactive user, got %d", w.Code)
	}

	var errResp dto.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp.Error.Code != "ACCOUNT_INACTIVE" {
		t.Errorf("expected ACCOUNT_INACTIVE code, got %s", errResp.Error.Code)
	}
}

func TestTeacherScopeMiddleware_ActiveTeacherSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uid := bson.NewObjectID()
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	ayID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjID := bson.NewObjectID()

	repo := &mockTeacherScopeRepo{
		userFunc: func(ctx context.Context, id bson.ObjectID) (*model.User, error) {
			return &model.User{
				ID:       uid,
				Role:     model.RoleTeacher,
				IsActive: true,
				SchoolID: schoolID,
			}, nil
		},
		teacherFunc: func(ctx context.Context, id bson.ObjectID) (*model.Teacher, error) {
			return &model.Teacher{
				ID:            teacherID,
				UserID:        uid,
				IsWalas:       true,
				WalasClassIDs: []bson.ObjectID{classID},
			}, nil
		},
		ayFunc: func(ctx context.Context, sID bson.ObjectID) (*model.AcademicYear, error) {
			return &model.AcademicYear{
				ID:       ayID,
				SchoolID: sID,
				Semester: 1,
				IsActive: true,
			}, nil
		},
		pairsFunc: func(ctx context.Context, sID, tID, aID bson.ObjectID) ([]repository.ClassSubjectPair, error) {
			return []repository.ClassSubjectPair{
				{ClassID: classID, SubjectID: subjID},
			}, nil
		},
	}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, uid.Hex())
		c.Next()
	})
	r.Use(TeacherScopeMiddleware(repo))
	r.GET("/test", func(c *gin.Context) {
		scope, ok := GetTeacherScope(c)
		if !ok || scope == nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if scope.TeacherID != teacherID || scope.AcademicYearID != ayID || !scope.IsWalas {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK for valid teacher scope, got %d", w.Code)
	}
}
