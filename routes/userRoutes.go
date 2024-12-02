package routes

import (
	"dsb/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func UserRoutes(r *gin.Engine) {
	// Serve Swagger UI (Swagger spec file can be found at /swagger.yaml)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define routes for user creation and login
	r.POST("/user", controllers.CreateUser)
	r.POST("/login", controllers.HandleLogin)
}
