package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Data:    data,
		Message: "success",
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    400,
		Data:    nil,
		Message: message,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    401,
		Data:    nil,
		Message: message,
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    403,
		Data:    nil,
		Message: message,
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    404,
		Data:    nil,
		Message: message,
	})
}

// InternalServerError 服务器内部错误
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    500,
		Data:    nil,
		Message: message,
	})
}