package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 需要重写
func Response(c *gin.Context, httpStatusCode int, code int, msg string, data interface{}) {
	c.JSON(httpStatusCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
func ResponseUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
func ResponseSuccess(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, SuccessCode, SuccessMsg, data)
}

func ResponseInternalError(c *gin.Context) {
	Response(c, http.StatusInternalServerError, InternalServerErrorCode, InternalServerErrorMsg, nil)
}

func ResponseError(c *gin.Context, code int, msg string) {
	Response(c, http.StatusInternalServerError, code, msg, nil)
}
