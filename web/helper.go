package web

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var isNum, isPassword *regexp.Regexp

func init() {
	isNum = regexp.MustCompile("^(\\d+)$")
	isPassword = regexp.MustCompile("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}$")
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
	if isPassword.MatchString(password) {
		return true
	}
	return false
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
