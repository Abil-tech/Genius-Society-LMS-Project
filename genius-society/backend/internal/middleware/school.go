package middleware

import "github.com/gin-gonic/gin"

// GetSchoolID extracts the school ID from gin context.
// It expects that a previous middleware (e.g., authentication) has stored the school ID under the "schoolID" key.
// If not set, it returns an empty string.
func GetSchoolID(c *gin.Context) string {
	if v, exists := c.Get("schoolID"); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
