package routes

import (
	"dsb/controllers"
	"dsb/middleware"

	"github.com/gin-gonic/gin"
)

func SuggestionRoutes(r *gin.Engine) {
	// Protected route group
	protected := r.Group("/")
	protected.Use(middleware.TokenAuthMiddleware())

	// POST /suggestion route is now protected
	protected.POST("/suggestion", controllers.PostSuggestion)
	protected.GET("/suggestion", controllers.GetAllSuggestions)
	protected.GET("/suggestion/user/:username", controllers.GetSuggestionsByUser)
	protected.GET("/suggestion/tag/:role", controllers.GetSuggestionWithTag)
}
