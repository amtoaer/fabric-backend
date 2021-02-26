package security

import "testing"

func TestSecurity(t *testing.T) {
	privateKey, publicKey, err := GenerateRsaKey()
	if err != nil {
		t.Error(err)
	}
	testMessage := "曲曲折折的荷塘上面，弥望的是田田的叶子。叶子出水很高，像亭亭的舞女的裙。"
	encrypt, _ := RsaEncrypt([]byte(testMessage), publicKey)
	decrypt, _ := RsaDecrypt(encrypt, privateKey)
	if string(decrypt) != testMessage {
		t.Error("经过加密->解密后数据出现变化")
	}
}
