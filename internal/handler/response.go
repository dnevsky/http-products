package handler

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func newErrorResponseJSON(c *gin.Context, statusCode int, message map[string]interface{}) {
	c.AbortWithStatusJSON(statusCode, message)
}
