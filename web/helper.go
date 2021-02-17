package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用来生成一个错误
func getError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": err,
	})
}
