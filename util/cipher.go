package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
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

func AesPkcs7Encrypt(key []byte, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return make([]byte, 0), err
	}

	src = PKCS7Padding(src, block.BlockSize())

	return ECBEncrypt(block, src, key)
}

func AesPkcs7Decrypt(key []byte, dst []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return make([]byte, 0), err
	}

	src, err := ECBDecrypt(block, dst, key)
	if err != nil {
		panic(err)
	}

	src = PKCS7UnPadding(src)

	return src, nil
}
