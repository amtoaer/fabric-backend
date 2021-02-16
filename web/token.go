package web

import (
	"time"

	"github.com/amtoaer/fabric-backend/model"
	jwt "github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("MYSECRETKEY")

// Claims token中存储的payload
type Claims struct {
	ID             uint
	StandardClaims jwt.StandardClaims
}

// Valid 实现要求的接口
func (c Claims) Valid() error { return nil }

// 通过用户主键生成token
func generateToken(ID uint, expireDuration time.Duration) (string, error) {
	// token过期时间
	expireTime := time.Now().Add(expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID: ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	})
	return token.SignedString(secretKey)
}

// 通过token得到用户信息
func parseToken(tokenStr string) (result *model.User, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok {
		result, err = model.GetUserByID(claims.ID)
	}
	return
}
