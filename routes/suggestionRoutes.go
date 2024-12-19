package routes

import (
	"dsb/controllers"
	"dsb/middleware"

	"github.com/gin-gonic/gin"
)

func SuggestionRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.TokenAuthMiddleware())
	protected.POST("/suggestion", controllers.PostSuggestion)
	protected.GET("/suggestion", controllers.GetAllSuggestions)
	protected.GET("/suggestion/user/:username", controllers.GetSuggestionsByUser)
	protected.GET("/suggestion/tag/:role", controllers.GetSuggestionWithTag)
	protected.POST("/suggestion/vote/:type", controllers.HandleVote)
	protected.GET("/admin/suggestions/:adminID", controllers.GetAdminSuggestions)
	protected.POST("/admin/suggestions/:suggestionId/respond", controllers.RespondToSuggestion)
	protected.POST("/suggestion/comment/:suggestionId", controllers.PostComment)
	protected.GET("/suggestions/count", controllers.GetUserSuggestions)
	protected.GET("/suggestion/:suggestionId/comments", controllers.GetComments)
}
