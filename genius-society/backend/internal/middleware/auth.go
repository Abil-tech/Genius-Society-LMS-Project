package middleware

import (
	"net/http"
	"strings"

	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID  = "userID"
	ContextKeyRole    = "userRole"
	ContextKeyIsWalas = "isWalas"
)

func Auth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try HTTP-only cookie first
		tokenStr, err := c.Cookie("token")
		if err != nil {
			// Fallback: Authorization: Bearer <token>
			authHeader := c.GetHeader("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   gin.H{"code": "UNAUTHORIZED", "message": "authentication required"},
				})
				return
			}
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims, err := authSvc.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   gin.H{"code": "UNAUTHORIZED", "message": "invalid or expired token"},
			})
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeyIsWalas, claims.IsWalas)
		c.Next()
	}
}

// RequireRoles checks that the authenticated user has one of the allowed roles.
func RequireRoles(roles ...model.Role) gin.HandlerFunc {
	allowed := make(map[model.Role]bool)
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get(ContextKeyRole)
		userRole, ok := role.(model.Role)
		if !ok || !allowed[userRole] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   gin.H{"code": "FORBIDDEN", "message": "access denied"},
			})
			return
		}
		c.Next()
	}
}

// PrincipalReadOnly blocks all non-GET requests for principal role.
func PrincipalReadOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ContextKeyRole)
		if role == model.RolePrincipal && c.Request.Method != http.MethodGet {
			// Exception: principal can mark their own notifications as read (PATCH)
			if strings.Contains(c.FullPath(), "/notifications/") && c.Request.Method == http.MethodPatch {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   gin.H{"code": "FORBIDDEN", "message": "principal role is read-only"},
			})
			return
		}
		c.Next()
	}
}

// GetUserID extracts userID from gin context.
func GetUserID(c *gin.Context) string {
	v, _ := c.Get(ContextKeyUserID)
	id, _ := v.(string)
	return id
}

// GetRole extracts role from gin context.
func GetRole(c *gin.Context) model.Role {
	v, _ := c.Get(ContextKeyRole)
	role, _ := v.(model.Role)
	return role
}

// GetIsWalas extracts isWalas flag from gin context.
func GetIsWalas(c *gin.Context) bool {
	v, _ := c.Get(ContextKeyIsWalas)
	isWalas, _ := v.(bool)
	return isWalas
}
