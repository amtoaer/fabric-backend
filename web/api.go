package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/amtoaer/fabric-backend/model"
	"github.com/amtoaer/fabric-backend/security"
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
			getError(c, err)
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
		getError(c, err)
		return
	}
	token, err := generateToken(user.ID, time.Hour*72)
	if err != nil {
		getError(c, err)
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
		getError(c, err)
		return
	}
	token, err := generateToken(user.ID, time.Hour*72)
	if err != nil {
		getError(c, err)
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
		getError(c, fmt.Errorf("只有医生才能添加病历"))
		return
	}
	patientName, patientIDNumber := c.Request.FormValue("patientName"), c.Request.FormValue("patientIDNumber")
	patient, err := model.SearchUser(patientIDNumber, patientName)
	if err != nil {
		getError(c, err)
		return
	}
	content := c.Request.FormValue("content")
	afterFirstEncrypt, err := security.RsaEncrypt([]byte(content), []byte(user.PublicKey))
	if err != nil {
		getError(c, err)
		return
	}
	afterSecondEncrypt, err := security.RsaEncrypt(afterFirstEncrypt, []byte(patient.PublicKey))
	if err != nil {
		getError(c, err)
		return
	}
	transactionID, err := helper.AddRecord(service.Record{
		ObjectType:  "recordObj",
		PatientID:   patientIDNumber,
		PatientName: patientName,
		DoctorID:    user.IDNumber,
		DoctorName:  user.Name,
		Content:     string(afterSecondEncrypt),
	})
	if err != nil {
		getError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"transactionID": transactionID,
	})
}
