package zCrypto

import (
	"crypto/aes"
	"crypto/rsa"
	"testing"
)

func TestMD5(t *testing.T) {
	result := MD5("test")
	expected := "098f6bcd4621d373cade4e832627b4f6"
	if result != expected {
		t.Errorf("MD5 failed: expected %s, got %s", expected, result)
	}
}

func TestSHA1(t *testing.T) {
	result := SHA1("test")
	expected := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	if result != expected {
		t.Errorf("SHA1 failed: expected %s, got %s", expected, result)
	}
}

func TestSHA256(t *testing.T) {
	result := SHA256("test")
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if result != expected {
		t.Errorf("SHA256 failed: expected %s, got %s", expected, result)
	}
}

func TestSHA512(t *testing.T) {
	result := SHA512("test")
	expected := "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"
	if result != expected {
		t.Errorf("SHA512 failed: expected %s, got %s", expected, result)
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	original := []byte("test data")

	// 测试Base64编码
	encoded := Base64Encode(original)
	expectedEncoded := "dGVzdCBkYXRh"
	if encoded != expectedEncoded {
		t.Errorf("Base64Encode failed: expected %s, got %s", expectedEncoded, encoded)
	}

	// 测试Base64解码
	decoded, err := Base64Decode(encoded)
	if err != nil {
		t.Errorf("Base64Decode failed: %v", err)
	}
	if string(decoded) != string(original) {
		t.Errorf("Base64Decode failed: expected %s, got %s", original, decoded)
	}
}

func TestBase64URLEncodeDecode(t *testing.T) {
	original := []byte("test data with special chars: /+=")

	// 测试URL安全的Base64编码
	encoded := Base64URLEncode(original)
	expectedEncoded := "dGVzdCBkYXRhIHdpdGggc3BlY2lhbCBjaGFyczogLys9"
	if encoded != expectedEncoded {
		t.Errorf("Base64URLEncode failed: expected %s, got %s", expectedEncoded, encoded)
	}

	// 测试URL安全的Base64解码
	decoded, err := Base64URLDecode(encoded)
	if err != nil {
		t.Errorf("Base64URLDecode failed: %v", err)
	}
	if string(decoded) != string(original) {
		t.Errorf("Base64URLDecode failed: expected %s, got %s", original, decoded)
	}
}

func TestHexEncodeDecode(t *testing.T) {
	original := []byte("test")

	// 测试十六进制编码
	encoded := HexEncode(original)
	expectedEncoded := "74657374"
	if encoded != expectedEncoded {
		t.Errorf("HexEncode failed: expected %s, got %s", expectedEncoded, encoded)
	}

	// 测试十六进制解码
	decoded, err := HexDecode(encoded)
	if err != nil {
		t.Errorf("HexDecode failed: %v", err)
	}
	if string(decoded) != string(original) {
		t.Errorf("HexDecode failed: expected %s, got %s", original, decoded)
	}
}

func TestAESCBC(t *testing.T) {
	// 测试数据
	original := []byte("test AES CBC encryption")
	key, _ := GenerateAESKey(AESKeySize16) // 128位密钥
	iv, _ := GenerateIV(aes.BlockSize)

	// 加密
	encrypted, err := AESEncrypt(original, key, iv, AESModeCBC)
	if err != nil {
		t.Errorf("AESEncrypt CBC failed: %v", err)
	}

	// 解密
	decrypted, err := AESDecrypt(encrypted, key, iv, AESModeCBC)
	if err != nil {
		t.Errorf("AESDecrypt CBC failed: %v", err)
	}

	// 验证结果
	if string(decrypted) != string(original) {
		t.Errorf("AES CBC failed: expected %s, got %s", original, decrypted)
	}
}

func TestAESGCM(t *testing.T) {
	// 测试数据
	original := []byte("test AES GCM encryption")
	key, _ := GenerateAESKey(AESKeySize24) // 192位密钥

	// GCM模式不需要手动提供IV，内部会生成
	encrypted, err := AESEncrypt(original, key, nil, AESModeGCM)
	if err != nil {
		t.Errorf("AESEncrypt GCM failed: %v", err)
	}

	// 解密
	decrypted, err := AESDecrypt(encrypted, key, nil, AESModeGCM)
	if err != nil {
		t.Errorf("AESDecrypt GCM failed: %v", err)
	}

	// 验证结果
	if string(decrypted) != string(original) {
		t.Errorf("AES GCM failed: expected %s, got %s", original, decrypted)
	}
}

func TestGenerateAESKey(t *testing.T) {
	// 测试生成16字节密钥
	key16, err := GenerateAESKey(AESKeySize16)
	if err != nil {
		t.Errorf("GenerateAESKey 16 failed: %v", err)
	}
	if len(key16) != AESKeySize16 {
		t.Errorf("GenerateAESKey 16: expected length %d, got %d", AESKeySize16, len(key16))
	}

	// 测试生成24字节密钥
	key24, err := GenerateAESKey(AESKeySize24)
	if err != nil {
		t.Errorf("GenerateAESKey 24 failed: %v", err)
	}
	if len(key24) != AESKeySize24 {
		t.Errorf("GenerateAESKey 24: expected length %d, got %d", AESKeySize24, len(key24))
	}

	// 测试生成32字节密钥
	key32, err := GenerateAESKey(AESKeySize32)
	if err != nil {
		t.Errorf("GenerateAESKey 32 failed: %v", err)
	}
	if len(key32) != AESKeySize32 {
		t.Errorf("GenerateAESKey 32: expected length %d, got %d", AESKeySize32, len(key32))
	}

	// 测试生成无效长度密钥
	_, err = GenerateAESKey(10)
	if err == nil {
		t.Error("GenerateAESKey should fail for invalid key size")
	}
}

func TestRSA(t *testing.T) {
	// 测试数据
	original := []byte("test RSA encryption")

	// 生成RSA密钥对
	privateKey, publicKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Errorf("GenerateRSAKeyPair failed: %v", err)
	}

	// 测试加密
	encrypted, err := RSAEncrypt(original, publicKey)
	if err != nil {
		t.Errorf("RSAEncrypt failed: %v", err)
	}

	// 测试解密
	decrypted, err := RSADecrypt(encrypted, privateKey)
	if err != nil {
		t.Errorf("RSADecrypt failed: %v", err)
	}

	// 验证结果
	if string(decrypted) != string(original) {
		t.Errorf("RSA failed: expected %s, got %s", original, decrypted)
	}

	// 测试签名和验证
	signature, err := RSASign(original, privateKey)
	if err != nil {
		t.Errorf("RSASign failed: %v", err)
	}

	err = RSAVerify(original, signature, publicKey)
	if err != nil {
		t.Errorf("RSAVerify failed: %v", err)
	}

	// 测试无效签名
	invalidSignature := append(signature[:len(signature)-1], signature[len(signature)-1]^0xFF)
	err = RSAVerify(original, invalidSignature, publicKey)
	if err == nil {
		t.Error("RSAVerify should fail for invalid signature")
	}
}

func TestRSAKeyPEMExportImport(t *testing.T) {
	// 生成RSA密钥对
	privateKey, publicKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Errorf("GenerateRSAKeyPair failed: %v", err)
	}

	// 测试导出私钥
	privatePEM, err := ExportRSAKeyToPEM(privateKey)
	if err != nil {
		t.Errorf("ExportRSAKeyToPEM private failed: %v", err)
	}

	// 测试导入私钥
	importedPrivate, err := ImportRSAKeyFromPEM(privatePEM)
	if err != nil {
		t.Errorf("ImportRSAKeyFromPEM private failed: %v", err)
	}

	// 验证私钥类型
	_, ok := importedPrivate.(*rsa.PrivateKey)
	if !ok {
		t.Error("ImportRSAKeyFromPEM private returned wrong type")
	}

	// 测试导出公钥
	publicPEM, err := ExportRSAKeyToPEM(publicKey)
	if err != nil {
		t.Errorf("ExportRSAKeyToPEM public failed: %v", err)
	}

	// 测试导入公钥
	importedPublic, err := ImportRSAKeyFromPEM(publicPEM)
	if err != nil {
		t.Errorf("ImportRSAKeyFromPEM public failed: %v", err)
	}

	// 验证公钥类型
	_, ok = importedPublic.(*rsa.PublicKey)
	if !ok {
		t.Error("ImportRSAKeyFromPEM public returned wrong type")
	}
}

func TestGenerateRandomString(t *testing.T) {
	// 测试生成不同长度的随机字符串
	lengths := []int{5, 10, 16, 32}
	for _, length := range lengths {
		str, err := GenerateRandomString(length)
		if err != nil {
			t.Errorf("GenerateRandomString(%d) failed: %v", length, err)
		}
		if len(str) != length {
			t.Errorf("GenerateRandomString(%d): expected length %d, got %d", length, length, len(str))
		}
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	// 测试生成不同长度的随机字节
	lengths := []int{5, 10, 16, 32}
	for _, length := range lengths {
		bytes, err := GenerateRandomBytes(length)
		if err != nil {
			t.Errorf("GenerateRandomBytes(%d) failed: %v", length, err)
		}
		if len(bytes) != length {
			t.Errorf("GenerateRandomBytes(%d): expected length %d, got %d", length, length, len(bytes))
		}
	}
}