package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/api"
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

	api.AddPhotosRoutes(v1)
	api.AddPhotosHomeRoutes(v1)
	api.AddDownloadRoutes(v1)
}

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		fmt.Println(c.Request.Method)

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
