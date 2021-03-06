package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/amtoaer/fabric-backend/model"
	"github.com/amtoaer/fabric-backend/security"
	"github.com/amtoaer/fabric-backend/service"
	"github.com/gin-gonic/gin"
)

type param struct {
	Name            string
	IDNumber        string
	Type            bool
	Password        string
	PatientIDNumber string
	DoctorIDNumber  string
	PublicKey       string
	PrivateKey      string
	Content         string
}

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
			getError(c, nil, "您访问的功能需要登录")
			c.Abort()
			return
		}
		user, err := parseToken(token)
		if err != nil {
			getError(c, err, "身份信息失效，请重新登录")
			c.Abort()
			return
		}
		c.Set("user", user)
	}
}

// 用户登录
// 请求属性 ID、Password
func login(c *gin.Context) {
	var params param
	if c.Bind(&params) != nil {
		getError(c, nil, "参数格式有误")
		return
	}
	IDNumber := params.IDNumber
	password := params.Password
	if !(checkIDNumber(IDNumber) && checkPassword(password)) {
		getError(c, nil, "参数内容有误")
		return
	}
	user, err := model.FindUser(IDNumber, password)
	if err != nil {
		getError(c, err, "账户不存在或密码错误，请重试")
		return
	}
	token, err := generateToken(user.ID, time.Hour*72)
	if err != nil {
		getError(c, err, "生成身份凭证失败，请重试")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"token":   token,
	})
}

// 用户注册
// 请求属性 IDNumber、Password、Name、Type
func register(c *gin.Context) {
	var params param
	if c.Bind(&params) != nil {
		getError(c, nil, "参数格式有误")
		return
	}
	IDNumber := params.IDNumber
	password := params.Password
	name := params.Name
	typ := params.Type
	if !(checkIDNumber(IDNumber) && checkPassword(password) && checkName(name)) {
		getError(c, nil, "参数内容有误")
		return
	}
	user, err := model.InsertUser(IDNumber, password, name, typ)
	if err != nil {
		getError(c, err, "注册用户失败，请重试")
		return
	}
	token, err := generateToken(user.ID, time.Hour*72)
	if err != nil {
		getError(c, err, "生成身份凭证失败，请重试")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"token":   token,
	})
}

// 添加病历
// 请求属性 PatientIDNumber、PublicKey、Content
func addRecord(c *gin.Context) {
	var params param
	if c.Bind(&params) != nil {
		getError(c, nil, "参数格式有误")
		return
	}
	tmp, _ := c.Get("user")
	user := tmp.(*model.User)
	if !user.Type {
		getError(c, nil, "只有医生才能添加病历")
		return
	}
	PatientIDNumber := params.PatientIDNumber
	if !checkIDNumber(PatientIDNumber) {
		getError(c, nil, "参数内容有误")
		return
	}
	patient, err := model.SearchUser(PatientIDNumber)
	if err != nil {
		getError(c, nil, "未找到病人信息")
		return
	}
	// 添加时需保证添加请求由医生发起，且病人信息存在
	if patient.Type == user.Type {
		getError(c, nil, "未找到病人信息")
		return
	}
	PublicKey := params.PublicKey
	if PublicKey != patient.PublicKey {
		getError(c, nil, "公钥内容不符合")
		return
	}
	Content := params.Content
	// 先用医生公钥加密，再用病人公钥加密
	afterFirstEncrypt, err := security.RsaEncrypt([]byte(Content), []byte(user.PublicKey))
	if err != nil {
		getError(c, err, "使用医生公钥加密信息失败")
		return
	}
	afterSecondEncrypt, err := security.RsaEncrypt(afterFirstEncrypt, []byte(patient.PublicKey))
	if err != nil {
		getError(c, err, "使用病人公钥加密信息失败")
		return
	}
	transactionID, err := helper.AddRecord(service.Record{
		ObjectType:     "recordObj",
		PatientID:      patient.IDNumber,
		PatientName:    patient.Name,
		DoctorID:       user.IDNumber,
		DoctorName:     user.Name,
		ContentEncrypt: afterSecondEncrypt,
	})
	if err != nil {
		getError(c, err, "添加病历失败，请重试")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"transactionID": transactionID,
	})
}

// 更新病历
// 请求属性 PatientIDNumber、PublicKey、Content
func updateRecord(c *gin.Context) {
	var params param
	if c.Bind(&params) != nil {
		getError(c, nil, "参数格式有误")
		return
	}
	tmp, _ := c.Get("user")
	user := tmp.(*model.User)
	if !user.Type {
		getError(c, nil, "只有医生才能修改病历")
		return
	}
	PatientIDNumber := params.PatientIDNumber
	if !checkIDNumber(PatientIDNumber) {
		getError(c, nil, "参数内容有误")
		return
	}
	patient, err := model.SearchUser(PatientIDNumber)
	if err != nil {
		getError(c, nil, "未找到病人信息")
		return
	}
	// 添加时需保证添加请求由医生发起，且病人信息存在
	if patient.Type == user.Type {
		getError(c, nil, "未找到病人信息")
		return
	}
	PublicKey := params.PublicKey
	if PublicKey != patient.PublicKey {
		getError(c, nil, "公钥内容不符合")
		return
	}
	Content := params.Content
	// 先用医生公钥加密，再用病人公钥加密
	afterFirstEncrypt, err := security.RsaEncrypt([]byte(Content), []byte(user.PublicKey))
	if err != nil {
		getError(c, err, "使用医生公钥加密信息失败")
		return
	}
	afterSecondEncrypt, err := security.RsaEncrypt(afterFirstEncrypt, []byte(patient.PublicKey))
	if err != nil {
		getError(c, err, "使用病人公钥加密信息失败")
		return
	}
	transactionID, err := helper.UpdateRecord(service.Record{
		ObjectType:     "recordObj",
		PatientID:      patient.IDNumber,
		PatientName:    patient.Name,
		DoctorID:       user.IDNumber,
		DoctorName:     user.Name,
		ContentEncrypt: afterSecondEncrypt,
	})
	if err != nil {
		getError(c, err, "更新病历失败，请重试")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"transactionID": transactionID,
	})
}

// 通过医生ID查询病历列表
// 请求属性 DoctorIDNumber
func searchRecordByDoctorID(c *gin.Context) {
	var params param
	if c.Bind(&params) != nil {
		getError(c, nil, "参数格式有误")
		return
	}
	DoctorIDNumber := params.DoctorIDNumber
	if !(checkIDNumber(DoctorIDNumber)) {
		getError(c, nil, "参数内容有误")
		return
	}
	result, err := helper.QueryRecordByDoctorID(DoctorIDNumber)
	if err != nil {
		getError(c, err, "查询失败，请重试")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// 通过病人IDNumber查询病历列表
// 请求属性 PatientIDNumber
func searchRecordByPatientID(c *gin.Context) {
	var params param
	if c.Bind(&params) != nil {
		getError(c, nil, "参数格式有误")
		return
	}
	PatientIDNumber := params.PatientIDNumber
	if !(checkIDNumber(PatientIDNumber)) {
		getError(c, nil, "参数内容有误")
		return
	}
	result, err := helper.QueryRecordByPatientID(PatientIDNumber)
	if err != nil {
		getError(c, err, "查询失败，请重试")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// 通过病人IDNumber和医生IDNumber及两者私钥得到病历详情
// 请求属性 IDNumber、PrivateKey
func searchRecordByKey(c *gin.Context) {
	var params param
	if c.Bind(&params) != nil {
		getError(c, nil, "参数格式有误")
		return
	}
	var doctorKey, patientKey string
	// 获取请求发起人
	tmp, _ := c.Get("user")
	firstUser := tmp.(*model.User)
	// 拿到请求参数中的IDNumber
	IDNumber := params.IDNumber
	if !checkIDNumber(IDNumber) {
		getError(c, nil, "参数内容有误")
		return
	}
	// 获取到另一个人的信息
	secondUser, err := model.SearchUser(IDNumber)
	if err != nil {
		getError(c, nil, "获取对方身份信息失败")
		return
	}
	PrivateKey := params.PrivateKey
	if PrivateKey != secondUser.PrivateKey {
		getError(c, nil, "私钥内容不符合")
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
		getError(c, err, "查询病历信息失败")
		return
	}
	var result service.Record
	err = json.Unmarshal([]byte(tmpResult), &result)
	if err != nil {
		getError(c, err, "解析病历信息失败")
		return
	}
	// 先用病人私钥解密，再用医生私钥解密
	afterFirstDecrypt, err := security.RsaDecrypt(result.ContentEncrypt, []byte(patientKey))
	if err != nil {
		getError(c, err, "使用病人私钥解密信息失败")
		return
	}
	afterSecondDecrypt, err := security.RsaDecrypt(afterFirstDecrypt, []byte(doctorKey))
	if err != nil {
		getError(c, err, "使用医生私钥解密信息失败")
		return
	}
	result.Content = string(afterSecondDecrypt)
	// 解密历史条目
	for index := range result.Historys {
		// 先用病人私钥解密，再用医生私钥解密
		afterFirstDecrypt, err := security.RsaDecrypt(result.Historys[index].History.ContentEncrypt, []byte(patientKey))
		if err != nil {
			getError(c, err, "使用病人私钥解密信息失败")
			return
		}
		afterSecondDecrypt, err := security.RsaDecrypt(afterFirstDecrypt, []byte(doctorKey))
		if err != nil {
			getError(c, err, "使用医生私钥解密信息失败")
			return
		}
		result.Historys[index].History.Content = string(afterSecondDecrypt)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}
