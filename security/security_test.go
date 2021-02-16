package security

import "testing"

func TestSecurity(t *testing.T) {
	privateKey, publicKey, err := GenerateRsaKey()
	if err != nil {
		t.Error(err)
	}
	testMessage := "这是一条测试信息"
	encrypt, _ := RsaEncrypt([]byte(testMessage), publicKey)
	decrypt, _ := RsaDecrypt(encrypt, privateKey)
	if string(decrypt) != testMessage {
		t.Error("经过加密->解密后数据出现变化")
	}
}
