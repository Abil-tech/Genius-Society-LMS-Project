package middleware

import (
	"context"
	"github.com/geniussociety/backend/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
)

// SetSchoolID middleware extracts the authenticated user's school ID from the database
// and stores it in the Gin context under the key "schoolID". Requests from users without
// a school are rejected.
func SetSchoolID(userRepo interface {
	FindByID(context.Context, bson.ObjectID) (*model.User, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user ID set by Auth middleware.
		userID := GetUserID(c)
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "user not authenticated"}})
			return
		}
		uid, err := bson.ObjectIDFromHex(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "BAD_REQUEST", "message": "invalid user id"}})
			return
		}
		// Load user to obtain SchoolID.
		usr, err := userRepo.FindByID(c.Request.Context(), uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": "failed to retrieve user"}})
			return
		}
		// Allow SuperAdmin without a school association.
		if role, exists := c.Get(ContextKeyRole); exists && role == model.RoleSuperAdmin {
			c.Next()
			return
		}
		if usr == nil || usr.SchoolID.IsZero() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": gin.H{"code": "FORBIDDEN", "message": "user does not belong to any school"}})
			return
		}
		// Store school ID as hex string.
		c.Set("schoolID", usr.SchoolID.Hex())
		c.Next()
	}
}
