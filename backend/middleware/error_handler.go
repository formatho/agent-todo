package middleware

import (
	"net/http"

	apperrors "github.com/formatho/agent-todo/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles AppError types
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's an AppError
			if appErr, ok := err.(*apperrors.AppError); ok {
				response := gin.H{
					"error": appErr.Message,
					"code":  appErr.Code,
				}

				if appErr.Details != nil {
					response["details"] = appErr.Details
				}

				c.JSON(appErr.HTTPStatus, response)
				c.Abort()
				return
			}

			// Generic error handling
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
	}
}

// HandleError is a helper function to handle errors in handlers
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		if appErr.Details != nil {
			response["details"] = appErr.Details
		}

		c.JSON(appErr.HTTPStatus, response)
		return
	}

	// Generic error
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
