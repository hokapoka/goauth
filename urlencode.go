package oauth

import (
	"bytes"
	"utf8"
	"fmt"
)

var safe = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz._-")

func EncodeURL(u string) string {
	var b bytes.Buffer
	for _, c := range u {
		b.Write(encodeChar(c))
	}
	return b.String()
}

func encodeChar(c int) []byte {
	b := make([]byte, utf8.RuneLen(c))
	utf8.EncodeRune(b, c)
	if bytes.Index(safe, b) != -1 {
		return b
	}
	return []byte(fmt.Sprintf("%%%02X", c))
}

