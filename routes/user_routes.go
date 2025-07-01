package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/controllers"
	"github.com/mahdi-cpp/PhotoKit/repositories"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo)

	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userController.ListUsers)
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
		userRoutes.PUT("/:id/online", userController.UpdateOnlineStatus)
	}
}
