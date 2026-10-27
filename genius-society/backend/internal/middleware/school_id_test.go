package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geniussociety/backend/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockUserRepoForMW struct {
	user *model.User
	err  error
}

func (m *mockUserRepoForMW) FindByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	return m.user, m.err
}

func TestSetSchoolID_Unauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockUserRepoForMW{}

	r := gin.New()
	// No auth context set
	r.Use(SetSchoolID(repo))
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing userID context, got %d", w.Code)
	}
}

func TestSetSchoolID_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockUserRepoForMW{}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, "invalid-hex-id")
		c.Next()
	})
	r.Use(SetSchoolID(repo))
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for invalid userID format, got %d", w.Code)
	}
}

func TestSetSchoolID_DBLookupError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := bson.NewObjectID()
	repo := &mockUserRepoForMW{
		err: errors.New("mongodb connection timeout"),
	}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, uid.Hex())
		c.Next()
	})
	r.Use(SetSchoolID(repo))
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error when DB lookup fails, got %d", w.Code)
	}
}

func TestSetSchoolID_UserNotFoundOrZeroSchoolID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := bson.NewObjectID()
	repo := &mockUserRepoForMW{
		user: nil, // user not found
	}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, uid.Hex())
		c.Next()
	})
	r.Use(SetSchoolID(repo))
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when user is nil or schoolID is zero, got %d", w.Code)
	}
}

func TestSetSchoolID_CannotBeOverriddenByClientHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := bson.NewObjectID()
	trustedSchoolID := bson.NewObjectID()
	repo := &mockUserRepoForMW{
		user: &model.User{
			ID:       uid,
			SchoolID: trustedSchoolID,
		},
	}

	var observedSchoolID string
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(ContextKeyUserID, uid.Hex())
		c.Next()
	})
	r.Use(SetSchoolID(repo))
	r.GET("/test", func(c *gin.Context) {
		observedSchoolID = GetSchoolID(c)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	// Spoofed client headers
	req.Header.Set("X-School-ID", bson.NewObjectID().Hex())
	req.Header.Set("School-ID", bson.NewObjectID().Hex())
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	if observedSchoolID != trustedSchoolID.Hex() {
		t.Fatalf("expected observed school ID to match trusted user DB record %s, got %s", trustedSchoolID.Hex(), observedSchoolID)
	}

}
