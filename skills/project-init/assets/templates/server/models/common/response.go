package common

import (
	"server/enum"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    enum.Code `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
}

func Success(c *gin.Context) {
	c.JSON(200, Response{
		Code:    enum.CodeSuccess,
		Message: enum.MsgSuccess,
		Data:    nil,
	})
}

// SuccessWithMsg 返回成功响应（带消息）
func SuccessWithMsg(c *gin.Context, message string, data any) {
	c.JSON(200, gin.H{
		"code":    enum.CodeSuccess,
		"message": message,
		"data":    data,
	})
}

func SuccessWithData(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code":    enum.CodeSuccess,
		"message": enum.MsgSuccess,
		"data":    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(resolveHTTPStatus(code), gin.H{
		"code":    code,
		"message": message,
	})
}

// Invalid 返回参数错误响应
func Invalid(c *gin.Context, message string) {
	c.JSON(400, gin.H{
		"code":    enum.CodeInvalidParameter,
		"message": message,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, message string) {
	c.JSON(500, gin.H{
		"code":    enum.CodeFailed,
		"message": message,
	})
}

func resolveHTTPStatus(code int) int {
	if code == enum.CodeInvalidParameter {
		return 400
	}
	if code == enum.CodeInsufficientPermissions {
		return 403
	}
	if code == enum.CodeTokenInvalid {
		return 401
	}
	if isAuthBusinessCode(code) {
		return 401
	}
	if code >= 500 {
		return 500
	}
	if code >= 400 {
		return 400
	}
	return 200
}

func isAuthBusinessCode(code int) bool {
	if code <= 0 {
		return false
	}
	codeStr := strconv.Itoa(code)
	return len(codeStr) >= 3 && codeStr[:3] == "401"
}
