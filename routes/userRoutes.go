package routes

import (
	"dsb/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/user", controllers.CreateUser)        // Define routes for user creation and login
	r.POST("/login", controllers.HandleLogin)      // Route for user login
	r.POST("/admin", controllers.RegisterAdmin)    // Route for admin registration
	r.POST("/admin/login", controllers.AdminLogin) // Route for admin login
}
