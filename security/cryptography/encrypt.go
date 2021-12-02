package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

/* Generate256BitRandomHexKey will generate a random 256 bit key encoded in base16 string */
func Generate256BitRandomHexKey() string {
	return EncodeHex(generate256BitRandomKey())
}

/* Generate256BitRandomBase64Key will generate a random 256 bit key encoded in base64 string */
func Generate256BitRandomBase64Key() string {
	return EncodeBASE64(generate256BitRandomKey())
}

func generate256BitRandomKey() []byte {
	bytes := make([]byte, 32)
	//generate a random 32 byte key
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	return bytes
}

// EncryptAES encrypt data []byte using hex encoded key using AES GCM. The key before encoded must be 32 bytes in length
func EncryptWithGCM(data []byte, key []byte) ([]byte, error) {
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data.
	//The first nonce argument in Seal is the prefix.
	cipherText := aesGCM.Seal(nonce, nonce, data, nil)
	return cipherText, nil
}

// DecryptAES encrypt data []byte using hex encoded key using AES GCM. The key before encoded must be 32 bytes in length
func DecryptWithGCM(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, cipherText := encryptedData[:aesGCM.NonceSize()], encryptedData[aesGCM.NonceSize():]

	return aesGCM.Open(nil, nonce, cipherText, nil)
}

func Encrypt(data []byte, byteKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(byteKey)
	if err != nil {
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize + len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	//returns to base64 encoded string
	return cipherText, nil
}

func Decrypt(encryptedByte []byte, byteKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(byteKey)
	if err != nil {
		return nil, err
	}

	if len(encryptedByte) < aes.BlockSize {
		return nil, fmt.Errorf("encryptedByte block size is less than %d bytes", aes.BlockSize)
	}

	decryptedByte := make([]byte, len(encryptedByte) - aes.BlockSize)

	iv := encryptedByte[:aes.BlockSize]
	encryptedByte = encryptedByte[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decryptedByte, encryptedByte)

	return decryptedByte, nil
}

