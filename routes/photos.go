package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/repositories"
	"net/http"
)

func AddPhotosRoutes(rg *gin.RouterGroup) {

	route := rg.Group("/photos")

	route.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, repositories.RestUser())
	})
}
