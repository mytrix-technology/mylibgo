package cryptography

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashSHA1(data []byte) string {
	h := sha1.New()
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}

func HashMD5(data []byte) string {
	h := md5.New()
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}

func HashUTPassword(password string) string {
	return HashSHA1([]byte(password))
}

func HashSHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HmacSHA256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}