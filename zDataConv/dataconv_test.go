package zDataConv

import (
	"fmt"
	"testing"
)

func TestConv(t *testing.T) {
	a, err := String2Int("10")
	fmt.Println("String2Int:", a, err)

	b := "test"
	fmt.Println("string", b)
	bEnc := Base64Encode(b)
	fmt.Println("bEnc", bEnc)
	bDec, err := Base64Decode(bEnc)
	fmt.Println("bDec", bDec)

	c := "http://www.abc.com"
	fmt.Println("string", c)
	cEnc := UrlEncode(c)
	fmt.Println("cEnc", cEnc)
	cDec, err := UrlDecode(cEnc)
	fmt.Println("cDec", cDec)
}
