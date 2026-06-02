// Package response 提供 HTTP 统一响应格式的辅助函数。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是统一的 HTTP 响应结构体。
type Response struct {
	// Code 业务状态码，0 表示成功，非 0 表示失败。
	Code int `json:"code"`
	// Message 响应消息描述。
	Message string `json:"message"`
	// Data 响应数据载荷，成功时为业务数据，失败时为 nil。
	Data interface{} `json:"data"`
} // @name Response

// Success 返回成功响应（HTTP 200，code=0）。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// Error 返回错误响应。
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
