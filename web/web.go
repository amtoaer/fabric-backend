package web

import (
	"github.com/gin-gonic/gin"
)

// NewRouter 返回路由
func NewRouter() *gin.Engine {
	router := gin.Default()
	userFunction := router.Group("/api/user")
	userFunction.POST("/login", login)
	userFunction.POST("/register", register)
	otherFunction := router.Group("/api/other", headerAuthorization())
	otherFunction.POST("/addRecord", addRecord)
	// TODO: 添加更多功能
	return router
}
