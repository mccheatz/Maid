package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func ECBEncrypt(block cipher.Block, src, key []byte) ([]byte, error) {
	blockSize := block.BlockSize()

	encryptData := make([]byte, 0)
	tmpData := make([]byte, blockSize)

	for index := 0; index < len(src); index += blockSize {
		block.Encrypt(tmpData, src[index:index+blockSize])
		encryptData = append(encryptData, tmpData...)
	}
	return encryptData, nil
}

func ECBDecrypt(block cipher.Block, src, key []byte) ([]byte, error) {
	dst := make([]byte, 0)

	blockSize := block.BlockSize()
	tmpData := make([]byte, blockSize)

	for index := 0; index < len(src); index += blockSize {
		block.Decrypt(tmpData, src[index:index+blockSize])
		dst = append(dst, tmpData...)
	}

	return dst, nil
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func AES_ECB_PKCS7Encrypt(key []byte, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	src = PKCS7Padding(src, block.BlockSize())

	return ECBEncrypt(block, src, key)
}

func AES_ECB_PKCS7Decrypt(key []byte, dst []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	src, err := ECBDecrypt(block, dst, key)
	if err != nil {
		return nil, err
	}

	src = PKCS7UnPadding(src)

	return src, nil
}

func AES_CBC_Encrypt(key []byte, data []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	encrypted := make([]byte, len(data))
	cbcEncrypter := cipher.NewCBCEncrypter(block, iv)
	cbcEncrypter.CryptBlocks(encrypted, data)

	return encrypted, nil
}

func AES_CBC_Decrypt(key []byte, data []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(data))
	cbcDecrypter := cipher.NewCBCDecrypter(block, iv)
	cbcDecrypter.CryptBlocks(decrypted, data)

	return decrypted, nil
}

func MD5Sum(data []byte) []byte {
	result := md5.Sum(data)
	return result[:]
}

func MD5Hex(input []byte) string {
	return hex.EncodeToString(MD5Sum(input))
}

func Xor(data []byte, key []byte) ([]byte, error) {
	if len(data) != len(key) {
		return nil, errors.New("data and key size do not match")
	}

	for i := 0; i < len(data); i++ {
		data[i] ^= key[i]
	}

	return data, nil
}
