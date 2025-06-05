package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/repository"
	"net/http"
)

func AddPhotosHomeRoutes(rg *gin.RouterGroup) {

	route := rg.Group("/photos")

	route.GET("/feed", func(context *gin.Context) {
		context.JSON(http.StatusOK, repository.RestFeed())
	})

	route.GET("/recent", func(context *gin.Context) {
		context.JSON(http.StatusOK, repository.RestRecentDays())
	})

	route.GET("/pinned", func(context *gin.Context) {
		context.JSON(http.StatusOK, repository.RestPinnedCollections())
	})

	route.GET("/trips", func(context *gin.Context) {
		context.JSON(http.StatusOK, repository.RestTrips())
	})

	route.GET("/albums", func(context *gin.Context) {
		context.JSON(http.StatusOK, repository.RestAlbums())
	})

	route.GET("/camera", func(context *gin.Context) {
		context.JSON(http.StatusOK, repository.RestCamera())
	})

	route.GET("/gallery", func(context *gin.Context) {
		context.JSON(http.StatusOK, repository.RestGallery())
	})

	route.GET("/subtitle", func(context *gin.Context) {
		repository.ReloadSubtitle()
		context.JSON(http.StatusOK, repository.RestFeed())
	})
}
