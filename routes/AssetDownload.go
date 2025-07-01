package routes

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/cache"
	"io"
	"net/http"
	"os"
	"strings"
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
			c.AbortWithStatusJSON(404, gin.H{"error": "File not found"})
			return
		}

		fileSize, err := getFileSize(filepath)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to get file size"})
			return
		}

		//etag, err := generateETag(filepath)
		//if err != nil {
		//	c.AbortWithStatusJSON(500, gin.H{"error": "Failed to generate ETag"})
		//	return
		//}

		//c.Header("Content-Type", "mage/jpeg")
		//c.Header("Content-Encoding", "identity") // Disable compression
		//c.Next()
		c.Header("Content-Length", fmt.Sprintf("%d", fileSize))
		//.Header("ETag", etag)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		c.Header("Accept-Ranges", "bytes")
		c.File(filepath)
	})
}

func getFileSize(filepath string) (int64, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func generateETag(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func apiDownloadThumb(route *gin.RouterGroup) {

	route.GET("/thumbnail/:filename", func(c *gin.Context) {

		filename := c.Param("filename")

		if strings.Contains(filename, "png") {
			imgData, exists := cache.GetIconCash(filename)
			if exists {
				c.Data(http.StatusOK, "image/png", imgData) // Adjust MIME type as necessary
			}
			return
		}

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
