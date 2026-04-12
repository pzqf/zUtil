package zRand

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"strings"
)

func init() {
}

// RandN return [0, n)
func RandN(n int) int {
	return rand.IntN(n)
}

// RandInterval return [min(a, b), max(a, b))
func RandInterval(a, b int) int {
	if a == b {
		return a
	}
	min, max := a, b
	if a > b {
		min, max = b, a
	}
	return min + RandN(max-min)
}

const letterString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandString 生成随机字符串
func RandString(n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(letterString[RandN(len(letterString))])
	}
	return sb.String()
}

// RandCustomString 使用自定义字符集生成随机字符串
func RandCustomString(n int, letter string) string {
	str := []byte(letter)
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(str[RandN(len(str))])
	}
	return sb.String()
}

// SecureInt 生成安全随机整数 [0, n)
func SecureInt(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("n must be positive")
	}
	var buf [4]byte
	if _, err := cryptorand.Read(buf[:]); err != nil {
		return 0, fmt.Errorf("crypto/rand read failed: %w", err)
	}
	v := int(binary.LittleEndian.Uint32(buf[:]))
	if v < 0 {
		v = -v
	}
	return v % n, nil
}

// SecureString 生成安全随机字符串
func SecureString(n int) (string, error) {
	result := make([]byte, n)
	if _, err := cryptorand.Read(result); err != nil {
		return "", fmt.Errorf("crypto/rand read failed: %w", err)
	}
	for i := range result {
		result[i] = letterString[int(result[i])%len(letterString)]
	}
	return string(result), nil
}

// SecureBytes 生成安全随机字节
func SecureBytes(n int) ([]byte, error) {
	result := make([]byte, n)
	if _, err := cryptorand.Read(result); err != nil {
		return nil, fmt.Errorf("crypto/rand read failed: %w", err)
	}
	return result, nil
}
