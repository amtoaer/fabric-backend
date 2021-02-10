package web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

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
