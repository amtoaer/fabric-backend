package web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

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

// 生成随机的公钥和私钥
func generateRsaKey() (prvKey []byte, pubKey []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvKey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKey = pem.EncodeToMemory(block)
	return
}

// 使用公钥对信息进行加密
func rsaEncrypt(data, keyBytes []byte) (result []byte, err error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return result, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub := pubInterface.(*rsa.PublicKey)
	result, err = rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	return
}

// 使用私钥对信息进行解密
func rsaDecrypt(data, keyBytes []byte) (result []byte, err error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return result, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	result, err = rsa.DecryptPKCS1v15(rand.Reader, priv, data)
	return
}

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
func parseToken(tokenStr string) (result *User, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok {
		result, err = getUserByID(claims.ID)
	}
	return
}
