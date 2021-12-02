package cryptography

import (
	"encoding/base64"
	"encoding/hex"
)

//EncodeBASE64 : Encode []byte to Base64 string.
func EncodeBASE64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

//DecodeBASE64 : Decrypt Base64. Input string, output string
func DecodeBASE64(text string) ([]byte, error) {
	byt, err := base64.StdEncoding.DecodeString(text)
	return byt, err
}

//EncodeBASE64URL : Encrypt to Base64URL. Input string, output text
func EncodeBASE64URL(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

//DecodeBASE64URL : Decrypt to Base64URL. Input string, output text
func DecodeBASE64URL(text string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(text)
}

func EncodeHex(data []byte) string {
	return hex.EncodeToString(data)
}

func DecodeHex(text string) ([]byte, error) {
	return hex.DecodeString(text)
}