package utils

import (
	"fmt"
	"testing"
)

func TestAESEncode(t *testing.T) {
	s, err := AESEncode("这个", "abc123")
	fmt.Println(err)
	fmt.Println(s)
}

func TestAESDecode(t *testing.T) {
	s, err := AESDecode("c9c1b4346b9959103d056b46a7999e67", "abc123")
	fmt.Println(err)
	fmt.Println(s)
}
