package routes

import (
	"dsb/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	// Define routes for user creation and login
	r.POST("/user", controllers.CreateUser)
	r.POST("/login", controllers.HandleLogin)
}
