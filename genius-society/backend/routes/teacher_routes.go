package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterTeacher registers all routes for the Teacher portal.
// As defined in project specifications, this entrypoint is wired in routes.Setup
// and will be populated across Modules 3.1 through 3.10.
func RegisterTeacher(protected *gin.RouterGroup) {
	// Root group for teacher endpoints
	// Sub-routes will be attached in respective modules (3.1 - 3.10)
	_ = protected
}
