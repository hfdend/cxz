package utils

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestTripleDESDecrypt(t *testing.T) {
	key, _ := hex.DecodeString("87D111C339265D7D87C8A81A88A096A4F138D4439425CBA5")
	data, _ := hex.DecodeString("dc31e046910c62e4f6f80761788781ad8e92f616499a00f8")
	//src := "5e8487e6"
	//key := "0123456789abcdef12343212"
	b, e := TripleDESDecrypt(data, key)
	fmt.Println(e)
	fmt.Println(string(b))
}

func TestTripleDESEncrypt(t *testing.T) {
	key, _ := hex.DecodeString("87D111C339265D7D87C8A81A88A096A4F138D4439425CBA5")
	data := "6228480469401649170"
	//b := hex.EncodeToString([]byte(data))
	s, err := TripleDESEncrypt([]byte(data), key)
	fmt.Println(err)
	fmt.Println()
	fmt.Println(fmt.Printf("%x\n", s))
}
