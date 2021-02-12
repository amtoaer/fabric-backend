package web

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 客户端用户
type User struct {
	// 主键ID
	ID   uint
	Name string
	// 身份证号
	IDNumber   string
	Password   string
	PrivateKey string
	PublicKey  string
}

var db *gorm.DB

// 初始化数据库连接池
func init() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// 通过ID获取User（验证token时用）
func getUserByID(ID uint) (*User, error) {
	var result *User
	status := db.First(result, ID)
	return result, status.Error
}

// 通过身份证号和密码获取User（登录时用）
func findUser(IDNumber, password string) (*User, error) {
	var result *User
	status := db.Where("IDNumber = ? AND password = ?", IDNumber, password).First(result)
	return result, status.Error
}

// 插入新的User（注册时用）
func insertUser(IDNumber, password, Name string) (err error) {
	var dataToInsert = &User{
		IDNumber: IDNumber,
		Password: password,
		Name:     Name,
	}
	privateKey, publicKey, err := generateRsaKey()
	if err != nil {
		return
	}
	dataToInsert.PrivateKey, dataToInsert.PublicKey = string(privateKey), string(publicKey)
	status := db.Create(dataToInsert)
	return status.Error
}
