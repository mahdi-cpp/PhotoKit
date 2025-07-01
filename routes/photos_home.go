package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/repositories"
	"net/http"
)

func AddPhotosHomeRoutes(rg *gin.RouterGroup) {

	route := rg.Group("/photos")

	//route.GET("/assets", func(context *gin.Context) {
	//	context.JSON(http.StatusOK, repositories.RestLibrary())
	//})
	//
	//route.GET("/library", func(context *gin.Context) {
	//	context.JSON(http.StatusOK, repositories.RestLibrary())
	//})
	route.GET("/collections", func(context *gin.Context) {
		context.JSON(http.StatusOK, repositories.RestCollections())
	})

	route.GET("/recent", func(context *gin.Context) {
		context.JSON(http.StatusOK, repositories.RestRecentDays())
	})

	route.GET("/pinned", func(context *gin.Context) {
		context.JSON(http.StatusOK, repositories.RestPinnedCollections())
	})

	route.GET("/trips", func(context *gin.Context) {
		context.JSON(http.StatusOK, repositories.RestTrips())
	})

	route.GET("/albums", func(context *gin.Context) {
		context.JSON(http.StatusOK, repositories.RestAlbums())
	})

	route.GET("/camera", func(context *gin.Context) {
		context.JSON(http.StatusOK, repositories.RestCamera())
	})

	//route.GET("/gallery", func(context *gin.Context) {
	//	context.JSON(http.StatusOK, repositories.RestGallery())
	//})
	//
	//route.GET("/subtitle", func(context *gin.Context) {
	//	repositories.ReloadSubtitle()
	//	context.JSON(http.StatusOK, repositories.RestLibrary())
	//})
}
