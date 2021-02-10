package web

import "testing"

func TestSecurity(t *testing.T) {
	privateKey, publicKey, err := generateRsaKey()
	if err != nil {
		t.Error(err)
	}
	testMessage := "这是一条测试信息"
	encrypt, _ := rsaEncrypt([]byte(testMessage), publicKey)
	decrypt, _ := rsaDecrypt(encrypt, privateKey)
	if string(decrypt) != testMessage {
		t.Error("经过加密->解密后数据出现变化")
	}
}
