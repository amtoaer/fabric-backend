package web

import (
	"net/http"
	"time"

	"github.com/amtoaer/fabric-backend/model"
	"github.com/amtoaer/fabric-backend/service"
	"github.com/gin-gonic/gin"
)

var helper service.Service

// SetService 设置web全局的service
func SetService(s service.Service) {
	helper = s
}

// 鉴权筛选
func headerAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		user, err := parseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err,
			})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

// 用户登录
func login(c *gin.Context) {
	ID := c.Request.FormValue("ID")
	password := c.Request.FormValue("Password")
	user, err := model.FindUser(ID, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	token, err := generateToken(user.ID, time.Hour*72)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"token":   token,
	})
}

// 用户注册
func register(c *gin.Context) {
	IDNumber := c.Request.FormValue("IDNumber")
	password := c.Request.FormValue("password")
	name := c.Request.FormValue("name")
	user, err := model.InsertUser(IDNumber, password, name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	token, err := generateToken(user.ID, time.Hour*72)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"token":   token,
	})
}

// 添加病历
func addRecord(c *gin.Context) {
	tmp, _ := c.Get("user")
	user := tmp.(*model.User)
	if !user.Type {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "只有医生才能添加病历",
		})
		return
	}
	patientName, patientIDNumber := c.Request.FormValue("patientName"), c.Request.FormValue("patientIDNumber")
	_, err := model.SearchUser(patientIDNumber, patientName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	content := c.Request.FormValue("content")
	// TODO: 对病历内容进行加密
	transactionID, err := helper.AddRecord(service.Record{
		ObjectType:  "test",
		PatientID:   patientIDNumber,
		PatientName: patientName,
		DoctorID:    user.IDNumber,
		DoctorName:  user.Name,
		Content:     content,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"transactionID": transactionID,
	})
}
