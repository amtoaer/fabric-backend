package model

import (
	"fmt"

	"github.com/amtoaer/fabric-backend/security"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 客户端用户
type User struct {
	ID         uint   // 主键编号
	Name       string // 姓名
	IDNumber   string //身份证号
	Password   string `json:"-"` // 密码
	Type       bool   // 身份标记（True表示医生，False表示病人）
	PrivateKey string //私钥
	PublicKey  string //公钥
}

var db *gorm.DB

// 初始化数据库连接池
func init() {
	var err error
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败，错误为：%v\n", err)
	}
}

// GetUserByID 通过ID获取User（验证token时用）
func GetUserByID(ID uint) (*User, error) {
	var result *User
	status := db.First(result, ID)
	return result, status.Error
}

// FindUser 通过用户编号和密码获取User（登录时用）
func FindUser(ID, password string) (*User, error) {
	var result *User
	status := db.Where("ID = ? AND Password = ?", ID, password).First(result)
	return result, status.Error
}

// SearchUser 通过身份证号获取User（查询用户是否存在时用）
func SearchUser(IDNumber string) (*User, error) {
	var result *User
	status := db.Where("IDNumber = ?", IDNumber).First(result)
	return result, status.Error
}

// InsertUser 插入新的User（注册时用）
func InsertUser(IDNumber, password, Name string, typ bool) (*User, error) {
	var dataToInsert = &User{
		IDNumber: IDNumber,
		Password: password,
		Name:     Name,
		Type:     typ,
	}
	privateKey, publicKey, err := security.GenerateRsaKey()
	if err != nil {
		return dataToInsert, err
	}
	dataToInsert.PrivateKey, dataToInsert.PublicKey = string(privateKey), string(publicKey)
	status := db.Create(dataToInsert)
	return dataToInsert, status.Error
}
