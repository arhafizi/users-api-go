package responses

import (
	"net/http"

	"example.com/api/internal/api/validation"
	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data any) {
	c.JSON(http.StatusCreated, BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string, err any) {
	var errors any
	switch v := err.(type) {
	case error:
		errors = v.Error()
	case *[]validation.ValidationError:
		errors = v
	case nil:
		errors = nil
	default:
		errors = "Invalid request"
	}

	c.JSON(http.StatusBadRequest, BaseResponse{
		Status:  "fail",
		Message: message,
		Errors:  errors,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, BaseResponse{
		Status:  "fail",
		Message: message,
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, BaseResponse{
		Status:  "error",
		Message: message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, BaseResponse{
		Status:  "fail",
		Message: message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, BaseResponse{
		Status:  "fail",
		Message: message,
	})
}

func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
