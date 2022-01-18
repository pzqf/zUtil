package contenttype


import (
	"fmt"
	"testing"
)

func TestContentType(t *testing.T) {
	fmt.Println(GetFileContentType("a.jpeg"))
}
