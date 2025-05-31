package dtos

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponseWithData defines the structure for a success response that includes data.
type SuccessResponseWithData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SendErrorResponse sends a standardized error response.
func SendErrorResponse(c *gin.Context, statusCode int, message string, details ...string) {
	errResponse := ErrorResponse{
		Error: message,
	}
	if len(details) > 0 {
		errResponse.Details = details[0]
	}
	c.AbortWithStatusJSON(statusCode, errResponse)
}

// SendSuccessResponse sends a standardized success response without data.
func SendSuccessResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, SuccessResponse{Message: message})
}

// SendSuccessResponseWithData sends a standardized success response with data.
func SendSuccessResponseWithData(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponseWithData{Message: message, Data: data})
}
