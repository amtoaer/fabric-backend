package web

import (
	"encoding/json"
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
		if token == "" {
			getError(c, fmt.Errorf("您访问的功能需要登录！"))
			return
		}
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
// 请求属性 ID、Password
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
// 请求属性 IDNumber、Password、Name
func register(c *gin.Context) {
	IDNumber := c.Request.FormValue("IDNumber")
	password := c.Request.FormValue("Password")
	name := c.Request.FormValue("Name")
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
// 请求属性 patientName、patientIDNumber、content
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
	// 添加时需保证添加请求由医生发起，且病人信息存在
	if patient.Type == user.Type {
		getError(c, fmt.Errorf("病人信息不符"))
		return
	}
	// BUG:需要先验证签名
	content := c.Request.FormValue("content")
	// 先用医生公钥加密，再用病人公钥加密
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

// 更新病历
func updateRecord(c *gin.Context) {}

// 通过医生ID查询病历列表
// 请求属性 doctorID
func searchRecordByDoctorID(c *gin.Context) {
	doctorID := c.Request.FormValue("doctorID")
	result, err := helper.QueryRecordByDoctorID(doctorID)
	if err != nil {
		getError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// 通过病人ID查询病历列表
// 请求属性 patientID
func searchRecordByPatientID(c *gin.Context) {
	patientID := c.Request.FormValue("patientID")
	result, err := helper.QueryRecordByPatientID(patientID)
	if err != nil {
		getError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// 通过病人ID和医生ID及两者私钥得到病历详情
// 请求属性 ID、privateKey、name
func searchRecordByKey(c *gin.Context) {
	var doctorKey, patientKey string
	// 获取请求发起人
	tmp, _ := c.Get("user")
	firstUser := tmp.(*model.User)
	// 拿到请求参数中的ID和name
	ID := c.Request.FormValue("ID")
	name := c.Request.FormValue("name")
	// 获取到另一个人的信息
	secondUser, err := model.SearchUser(ID, name)
	if err != nil {
		getError(c, err)
		return
	}
	privateKey := c.Request.FormValue("privateKey")
	if privateKey != secondUser.PrivateKey {
		getError(c, fmt.Errorf("私钥内容不符合"))
		return
	}
	var tmpResult string
	// 发起人是医生
	if firstUser.Type {
		tmpResult, err = helper.QueryRecordByKey(secondUser.IDNumber, firstUser.IDNumber)
		doctorKey, patientKey = firstUser.PrivateKey, secondUser.PrivateKey
	} else {
		tmpResult, err = helper.QueryRecordByKey(firstUser.IDNumber, secondUser.IDNumber)
		patientKey, doctorKey = firstUser.PrivateKey, secondUser.PrivateKey
	}
	if err != nil {
		getError(c, err)
		return
	}
	var result *service.Record
	err = json.Unmarshal([]byte(tmpResult), result)
	if err != nil {
		getError(c, err)
		return
	}
	// 先用病人私钥解密，再用医生私钥解密
	afterFirstDecrypt, err := security.RsaDecrypt([]byte(result.Content), []byte(patientKey))
	if err != nil {
		getError(c, err)
		return
	}
	afterSecondDecrypt, err := security.RsaDecrypt(afterFirstDecrypt, []byte(doctorKey))
	if err != nil {
		getError(c, err)
		return
	}
	result.Content = string(afterSecondDecrypt)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}
