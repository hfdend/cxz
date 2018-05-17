package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

var (
	AesKey   = "astaxie12798akljzmknm.ahkjkljl;k"
	commonIV = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	block, _ = aes.NewCipher([]byte(AesKey))
)

func AesEncode(v string) string {
	plaintext := []byte(v)
	cfb := cipher.NewCFBEncrypter(block, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	result := fmt.Sprintf("%x", ciphertext)
	return result
}

func AesDecode(s string) string {
	ciphertext, _ := hex.DecodeString(s)
	cfbdec := cipher.NewCFBDecrypter(block, commonIV)
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	result := fmt.Sprintf("%s", plaintextCopy)
	return result
}
