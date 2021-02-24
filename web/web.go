package web

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// NewRouter 返回路由
func NewRouter() *gin.Engine {
	staticFS := getFileSystem()
	router := gin.Default()
	router.Use(static.Serve("/", staticFS))
	router.NoRoute(func(c *gin.Context) {
		c.FileFromFS("/index.html", staticFS)
	})
	// 用户功能
	user := router.Group("/api/user")
	user.POST("/login", login)
	user.POST("/register", register)
	// 搜索功能
	search := router.Group("/api/search")
	search.POST("/byDoctorID", searchRecordByDoctorID)
	search.POST("byPatientID", searchRecordByPatientID)
	search.POST("byKey", searchRecordByKey)
	// 修改功能
	modify := router.Group("/api/modify")
	modify.Use(headerAuthorization())
	modify.POST("/addRecord", addRecord)
	modify.POST("/updateRecord", updateRecord)
	return router
}
