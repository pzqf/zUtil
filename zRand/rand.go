package zRand

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandN return [0 , n)
func RandN(n int) int {
	return rand.Intn(n)
}

//RandInterval return [ min(a, b), max(a, b) )
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

func RandString(n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += fmt.Sprintf("%c", letterString[RandN(len(letterString))])
	}
	return res
}

func RandCustomString(n int, letter string) string {
	str := []byte(letter)
	res := ""
	for i := 0; i < n; i++ {
		res += fmt.Sprintf("%c", str[RandN(strings.Count(letter, "")-1)])
	}
	return res
}
