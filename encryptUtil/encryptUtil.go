package encryptUtil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io"
)

const (
	INTERNAL_KEY    = "5cb602c24a3211e6a797567a637f5735"
	CONFIG_ITEM_KEY = "auth.encrypt.key"
)

// Encryption using AES method, then encode using base64
func AesEncrypt(key, plainText string) (encryptTextStr string, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return "", err
	}
	plainByte := []byte(plainText)
	encryptByte := make([]byte, aes.BlockSize+len(plainText))
	// initialization vector
	iv := encryptByte[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(encryptByte[aes.BlockSize:], plainByte)
	encryptTextStr = string(Base64Encode(encryptByte))
	return
}

// Decryption using AES method
func AesDecrypt(key, encryptText string) (plainText string, err error) {
	var block cipher.Block
	var encryptByte []byte
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return
	}
	if encryptByte, err = Base64Decode([]byte(encryptText)); err != nil {
		return
	}
	//encryptByte := []byte(encryptText)
	if len(encryptByte) < aes.BlockSize {
		err = errors.New("encryptText too short")
		return
	}

	iv := encryptByte[:aes.BlockSize]
	encryptByte = encryptByte[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(encryptByte, encryptByte)

	plainText = string(encryptByte)
	return
}

// encode encodes a value using base64.
func Base64Encode(value []byte) []byte {
	//value := []byte(valueString)
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

// decode decodes a cookie using base64.
func Base64Decode(value []byte) ([]byte, error) {
	//value := []byte(valueString)
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

// Get Configured Encrypt/Decrypt Key
func GetConfiguredKey() string {
	return beego.AppConfig.String(CONFIG_ITEM_KEY)
}

// Calculate using MD5 method
func MD5Sum(value string) string {
	md5Sum := md5.Sum([]byte(value))
	return fmt.Sprintf("%x", md5Sum)
}
