package zDataConv

import (
	"encoding/base64"
	"strconv"
)

func String2Byte(str string) []byte {
	return []byte(str)
}

func String2Int(str string) (int, error) {
	return strconv.Atoi(str)
}

func String2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func String2Float32(str string) (float32, error) {
	f64, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, err
	}
	return float32(f64), nil
}

func String2Float64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func Byte2String(data []byte) string {
	return string(data)
}

func Int2String(data int) string {
	return strconv.Itoa(data)
}

func Int642String(data int64) string {
	return strconv.FormatInt(data, 10)
}

func Float322String(data float32) string {
	return strconv.FormatFloat(float64(data), 'g', -1, 32)
}

func Float642String(data float64) string {
	return strconv.FormatFloat(data, 'g', -1, 64)
}

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(sDec), nil
}

func UrlEncode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func UrlDecode(data string) (string, error) {
	sDec, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(sDec), nil
}
