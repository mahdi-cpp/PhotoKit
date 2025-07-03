package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/routes"
)

var (
	router = gin.Default()
)

func Run() {

	router.Use(CORSMiddleware())

	getRoutes()

	err := router.Run(":8095")
	if err != nil {
		fmt.Println("Error] failed to start Gin server due to: ", err.Error())
		return
	}
}

func getRoutes() {

	v1 := router.Group("/v1")

	routes.AddPhotosRoutes(v1)
	routes.AddPhotosHomeRoutes(v1)
	routes.AddDownloadRoutes(v1)
}

func CORSMiddleware() gin.HandlerFunc {

	allowedOrigins := map[string]bool{
		"https://yourdomain.com": true,
		"http://localhost:3000":  true,
	}

	return func(c *gin.Context) {

		if allowedOrigins[c.Request.Header.Get("Origin")] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			c.Writer.Header().Set("Vary", "Origin") // Important for caching
		}

		// Allow specific origins (replace "*" in production!)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Cache preflight response for 1 day (optional)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Allowed HTTP methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allowed headers (simplified + security-critical ones)
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Origin, Content-Type, Accept, Authorization, X-Requested-With")

		// Headers exposed to the client
		c.Writer.Header().Set("Access-Control-Expose-Headers",
			"Content-Length, Content-Type, ETag")

		// Allow credentials (if needed, but avoid with "*" origin)
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle OPTIONS preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 204 No Content is more correct for OPTIONS
			return
		}

		c.Next()

	}
}
