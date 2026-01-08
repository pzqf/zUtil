package zCrypto

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash"
	"io"
)

// MD5 计算字符串的MD5哈希值
func MD5(data string) string {
	return Hash(data, md5.New())
}

// SHA1 计算字符串的SHA1哈希值
func SHA1(data string) string {
	return Hash(data, sha1.New())
}

// SHA256 计算字符串的SHA256哈希值
func SHA256(data string) string {
	return Hash(data, sha256.New())
}

// SHA512 计算字符串的SHA512哈希值
func SHA512(data string) string {
	return Hash(data, sha512.New())
}

// Hash 通用哈希计算函数
func Hash(data string, h hash.Hash) string {
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Base64Encode Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// Base64URLEncode URL安全的Base64编码
func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64URLDecode URL安全的Base64解码
func Base64URLDecode(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

// HexEncode 十六进制编码
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode 十六进制解码
func HexDecode(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

// AES相关常量
const (
	AESKeySize16 = 16 // 128位
	AESKeySize24 = 24 // 192位
	AESKeySize32 = 32 // 256位
)

// AES加密模式类型
type AESMode string

const (
	AESModeCBC AESMode = "CBC"
	AESModeGCM AESMode = "GCM"
)

// AESEncrypt AES加密
func AESEncrypt(data, key, iv []byte, mode AESMode) ([]byte, error) {
	if len(key) != AESKeySize16 && len(key) != AESKeySize24 && len(key) != AESKeySize32 {
		return nil, fmt.Errorf("invalid AES key size: must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case AESModeCBC:
		if len(iv) != aes.BlockSize {
			return nil, fmt.Errorf("invalid IV size for CBC mode: must be %d bytes", aes.BlockSize)
		}
		return aesCBCEncrypt(block, data, iv)
	case AESModeGCM:
		return aesGCMEncrypt(block, data, iv)
	default:
		return nil, fmt.Errorf("unsupported AES mode: %s", mode)
	}
}

// AESDecrypt AES解密
func AESDecrypt(data, key, iv []byte, mode AESMode) ([]byte, error) {
	if len(key) != AESKeySize16 && len(key) != AESKeySize24 && len(key) != AESKeySize32 {
		return nil, fmt.Errorf("invalid AES key size: must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case AESModeCBC:
		if len(iv) != aes.BlockSize {
			return nil, fmt.Errorf("invalid IV size for CBC mode: must be %d bytes", aes.BlockSize)
		}
		return aesCBCDecrypt(block, data, iv)
	case AESModeGCM:
		return aesGCMDecrypt(block, data, iv)
	default:
		return nil, fmt.Errorf("unsupported AES mode: %s", mode)
	}
}

// aesCBCEncrypt 使用CBC模式加密
func aesCBCEncrypt(block cipher.Block, data, iv []byte) ([]byte, error) {
	// 填充数据
	data = pkcs7Padding(data, block.BlockSize())

	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	mode.CryptBlocks(encrypted, data)

	return encrypted, nil
}

// aesCBCDecrypt 使用CBC模式解密
func aesCBCDecrypt(block cipher.Block, data, iv []byte) ([]byte, error) {
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	// 去填充
	decrypted, err := pkcs7Unpadding(decrypted, block.BlockSize())
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// aesGCMEncrypt 使用GCM模式加密
func aesGCMEncrypt(block cipher.Block, data, iv []byte) ([]byte, error) {
	// 如果未提供IV，生成一个随机IV
	if iv == nil {
		iv = make([]byte, 12) // GCM模式推荐使用12字节的IV
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
	} else if len(iv) != 12 {
		return nil, fmt.Errorf("invalid IV size for GCM mode: must be 12 bytes")
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 加密并附加IV
	encrypted := aead.Seal(nil, iv, data, nil)
	return append(iv, encrypted...), nil
}

// aesGCMDecrypt 使用GCM模式解密
func aesGCMDecrypt(block cipher.Block, data, iv []byte) ([]byte, error) {
	// GCM模式使用12字节IV
	const gcmIVLength = 12

	// 如果未提供IV，从数据开头提取
	if iv == nil {
		if len(data) < gcmIVLength {
			return nil, fmt.Errorf("data too short to extract IV")
		}
		iv = data[:gcmIVLength]
		data = data[gcmIVLength:]
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, iv, data, nil)
}

// pkcs7Padding PKCS#7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpadding PKCS#7去填充
func pkcs7Unpadding(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	
	padding := int(data[length-1])
	if padding < 1 || padding > blockSize {
		return nil, fmt.Errorf("invalid padding")
	}

	return data[:(length - padding)], nil
}

// GenerateAESKey 生成指定长度的AES密钥
func GenerateAESKey(keySize int) ([]byte, error) {
	if keySize != AESKeySize16 && keySize != AESKeySize24 && keySize != AESKeySize32 {
		return nil, fmt.Errorf("invalid AES key size: must be 16, 24, or 32 bytes")
	}

	key := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateIV 生成指定长度的IV
func GenerateIV(size int) ([]byte, error) {
	iv := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	return iv, nil
}

// GenerateRSAKeyPair 生成RSA密钥对
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// RSAEncrypt RSA加密
func RSAEncrypt(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

// RSADecrypt RSA解密
func RSADecrypt(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
}

// RSASign RSA签名
func RSASign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}

// RSAVerify RSA签名验证
func RSAVerify(data, signature []byte, publicKey *rsa.PublicKey) error {
	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}

// ExportRSAKeyToPEM 将RSA密钥导出为PEM格式
func ExportRSAKeyToPEM(key interface{}) (string, error) {
	var block *pem.Block

	switch k := key.(type) {
	case *rsa.PrivateKey:
		block = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}
	case *rsa.PublicKey:
		bytes, err := x509.MarshalPKIXPublicKey(k)
		if err != nil {
			return "", err
		}
		block = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytes,
		}
	default:
		return "", fmt.Errorf("unsupported key type")
	}

	return string(pem.EncodeToMemory(block)), nil
}

// ImportRSAKeyFromPEM 从PEM格式导入RSA密钥
func ImportRSAKeyFromPEM(pemData string) (interface{}, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "RSA PUBLIC KEY":
		return x509.ParsePKIXPublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		if _, err := io.ReadFull(rand.Reader, result[i:i+1]); err != nil {
			return "", err
		}
		result[i] = charset[int(result[i])%len(charset)]
	}
	return string(result), nil
}

// GenerateRandomBytes 生成指定长度的随机字节
func GenerateRandomBytes(length int) ([]byte, error) {
	result := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, result)
	return result, err
}