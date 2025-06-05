package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/cache"
	"net/http"
)

func AddDownloadRoutes(rg *gin.RouterGroup) {
	route := rg.Group("/download")
	apiOriginalDownload(route)
	apiDownloadThumb(route)
	apiIcon(route)
}

func apiOriginalDownload(route *gin.RouterGroup) {

	route.GET("/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filepath, err := cache.SearchFile(filename)
		if err != nil {
			return
		}
		c.File(filepath)
	})
}

func apiDownloadThumb(route *gin.RouterGroup) {

	route.GET("/thumbnail/:filename", func(c *gin.Context) {

		filename := c.Param("filename")

		imgData, exists := cache.GetThumbCash(filename)
		if exists {
			fmt.Println("RAM")
			c.Data(http.StatusOK, "image/jpeg", imgData) // Adjust MIME type as necessary
		} else {

			filepath, err := cache.SearchFile(filename)
			if err != nil {
				fmt.Println("SearchFile error", err)
				return
			}

			fmt.Println("SSD")
			c.File(filepath)
			cache.AddThumbCash(filepath, filename)
		}
	})
}

func apiIcon(route *gin.RouterGroup) {

	route.GET("/icons/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		imgData, exists := cache.GetIconCash(filename)
		if exists {
			c.Data(http.StatusOK, "image/png", imgData) // Adjust MIME type as necessary
		}
	})
}
