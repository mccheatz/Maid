package util

import (
	"bytes"
	"fmt"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int, runes_optional ...[]rune) string {
	var runes = letterRunes
	if len(runes_optional) != 0 {
		runes = runes_optional[0]
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func RandChinese(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rune(rand.Intn(0x9eff-0x4e00) + 0x4e00)
	}
	return string(b)
}

func RandMacAddress() string {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	// Set the local bit
	buf[0] |= 2

	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

func ToBinaryString(data []byte) string {
	var buf bytes.Buffer
	for i := 0; i < len(data); i++ {
		for j := 0; j < 8; j++ {
			buf.WriteByte('0' + (data[i] >> (7 - j) & 1))
		}
	}
	return buf.String()
}
