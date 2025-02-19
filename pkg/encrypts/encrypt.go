package encrypts

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

var Key = "1111111111111111"

func Encrypt(plainText string, key string) (string, error) {
	keyBytes := []byte(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	plainBytes := []byte(plainText)

	// 添加填充到明文
	padding := aes.BlockSize - len(plainBytes)%aes.BlockSize
	paddedPlainBytes := make([]byte, len(plainBytes)+padding)
	copy(paddedPlainBytes, plainBytes)
	for i := len(plainBytes); i < len(paddedPlainBytes); i++ {
		paddedPlainBytes[i] = byte(padding)
	}

	ciphertext := make([]byte, len(paddedPlainBytes))

	// 使用CBC模式加密
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(ciphertext, paddedPlainBytes)

	cipherHex := hex.EncodeToString(ciphertext)

	return cipherHex, nil
}

func Decrypt(cipherHex string, key string) (string, error) {
	keyBytes := []byte(key)

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 使用CBC模式解密
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCBCDecrypter(block, iv)

	paddedPlainBytes := make([]byte, len(ciphertext))
	stream.CryptBlocks(paddedPlainBytes, ciphertext)

	// 移除填充字符
	padding := paddedPlainBytes[len(paddedPlainBytes)-1]
	unpaddedPlainBytes := paddedPlainBytes[:len(paddedPlainBytes)-int(padding)]

	return string(unpaddedPlainBytes), nil
}

func removePadding(data []byte) []byte {
	padding := data[len(data)-1]
	return data[:len(data)-int(padding)]
}
