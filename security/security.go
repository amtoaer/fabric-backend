package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// GenerateRsaKey 生成随机的公钥和私钥
func GenerateRsaKey() (prvKey []byte, pubKey []byte, err error) {
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

// RsaEncrypt 使用公钥对信息进行分段加密
func RsaEncrypt(data, keyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	var (
		dataSize = len(data)
		offSet   = 0
		once     = pub.Size() - 11
		endIndex int
		buffer   bytes.Buffer
	)
	for offSet < dataSize {
		endIndex = offSet + once
		if endIndex > dataSize {
			endIndex = dataSize
		}
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}

// RsaDecrypt 使用私钥对信息进行解密
func RsaDecrypt(data, keyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	var (
		dataSize = len(data)
		once     = priv.Size()
		offSet   = 0
		endIndex int
		buffer   bytes.Buffer
	)
	for offSet < dataSize {
		endIndex = offSet + once
		if endIndex > dataSize {
			endIndex = dataSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, priv, data[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}
