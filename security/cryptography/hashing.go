/*
Taken from andrewtooyut common crypto package
 */
package cryptography

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
)

// EncodeSHA1HMACBase64 : encrypt to SHA1HMAC input key, data String. Output to String in Base64 format
func EncodeSHA1HMACBase64(key string, data ...string) string {
	return EncodeBASE64(ComputeSHA1HMAC(key, data...))
}

// EncodeSHA1HMAC : encrypt to SHA1HMAC input key, data String. Output to String in Base16/Hex format
func EncodeSHA1HMAC(key string, data ...string) string {
	return fmt.Sprintf("%x", ComputeSHA1HMAC(key, data...))
}

//ComputeSHA1HMAC : encrypt to SHA1HMAC input key, data String. Output to String
func ComputeSHA1HMAC(key string, data ...string) []byte {
	h := hmac.New(sha1.New, []byte(key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return h.Sum(nil)
}

func EncodeSHA256HMACBase64(key string, data ...string) string {
	return EncodeBASE64(ComputeSHA256HMAC(key, data...))
}

func EncodeSHA256HMAC(key string, data ...string) string {
	return fmt.Sprintf("%x", ComputeSHA256HMAC(key, data...))
}

func ComputeSHA256HMAC(key string, data ...string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return h.Sum(nil)
}

func EncodeSHA512HMACBase64(key string, data ...string) string {
	return EncodeBASE64(ComputeSHA512HMAC(key, data...))
}

func EncodeSHA512HMAC(key string, data ...string) string {
	return fmt.Sprintf("%x", ComputeSHA512HMAC(key, data...))
}

func ComputeSHA512HMAC(key string, data ...string) []byte {
	h := hmac.New(sha512.New, []byte(key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return h.Sum(nil)
}
