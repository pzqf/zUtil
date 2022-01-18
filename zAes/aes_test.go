package zAes

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"
)

func TestAes(t *testing.T) {
	origData := []byte("无名科技")        // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	log.Println("原文：", string(origData))

	log.Println("------------------ CBC模式 --------------------")
	encrypted := EncryptCBC(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := DecryptCBC(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	log.Println("------------------ ECB模式 --------------------")
	encrypted = EncryptECB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = DecryptECB(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	log.Println("------------------ CFB模式 --------------------")
	encrypted = EncryptCFB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = DecryptCFB(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}
