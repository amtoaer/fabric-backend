package web

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var isNum *regexp.Regexp

func init() {
	isNum = regexp.MustCompile("^(\\d+)$")
}

// 用来生成一个错误
func getError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": err.Error(),
	})
}

func checkID(ID string) bool {
	if len(ID) > 0 && isNum.MatchString(ID) {
		return true
	}
	return false
}

func checkPassword(password string) bool {
	if len(password) == 0 || len(password) > 16 {
		return false
	}
	for i := 0; i < len(password); i++ {
		if !((password[i] >= '0' && password[i] <= '9') || (password[i] >= 'a' && password[i] <= 'z') || (password[i] >= 'A' && password[i] <= 'Z')) {
			return false
		}
	}
	return true
}

func checkName(name string) bool {
	if len(name) > 0 && len(name) <= 10 {
		return true
	}
	return false
}

func checkIDNumber(IDNumber string) bool {
	if len(IDNumber) == 18 && isNum.MatchString(IDNumber) {
		return true
	}
	return false
}
