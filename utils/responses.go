package utils

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Message string `json:"message"`
}

func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{Message: message})
}

func SendSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}
