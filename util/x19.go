package util

import (
	"errors"
	"math"
	"math/rand"
)

func X19PickKey(query byte) []byte {
	keys := []string{
		"MK6mipwmOUedplb6",
		"OtEylfId6dyhrfdn",
		"VNbhn5mvUaQaeOo9",
		"bIEoQGQYjKd02U0J",
		"fuaJrPwaH2cfXXLP",
		"LEkdyiroouKQ4XN1",
		"jM1h27H4UROu427W",
		"DhReQada7gZybTDk",
		"ZGXfpSTYUvcdKqdY",
		"AZwKf7MWZrJpGR5W",
		"amuvbcHw38TcSyPU",
		"SI4QotspbjhyFdT0",
		"VP4dhjKnDGlSJtbB",
		"UXDZx4KhZywQ2tcn",
		"NIK73ZNvNqzva4kd",
		"WeiW7qU766Q1YQZI",
	}
	return []byte(keys[query>>4&0xf])
}

func X19HttpEncrypt(bodyIn []byte) ([]byte, error) {
	body := make([]byte, int(math.Ceil(float64(len(bodyIn)+16)/16))*16)
	copy(body, bodyIn)
	randFill := []byte(RandStringRunes(0x10))
	for i := 0; i < len(randFill); i++ {
		body[i+len(bodyIn)] = randFill[i]
	}

	keyQuery := byte(rand.Intn(15))<<4 | 2
	initVector := []byte(RandStringRunes(0x10))
	encrypted, err := AES_CBC_Encrypt(X19PickKey(keyQuery), body, initVector)
	if err != nil {
		return nil, err
	}

	result := make([]byte, 16 /* iv */ +len(encrypted) /* encrypted (body + scissor) */ +1 /* key query */)
	for i := 0; i < 16; i++ {
		result[i] = initVector[i]
	}
	for i := 0; i < len(encrypted); i++ {
		result[i+16] = encrypted[i]
	}

	result[len(result)-1] = keyQuery

	return result, nil
}

func X19HttpDecrypt(body []byte) ([]byte, error) {
	if len(body) < 0x12 {
		return nil, errors.New("input body too short")
	}

	result, err := AES_CBC_Decrypt(X19PickKey(body[len(body)-1]), body[16:len(body)-1], body[:16])
	if err != nil {
		return nil, err
	}

	scissor := 0
	scissorPos := len(result) - 1
	for scissor < 16 {
		if result[scissorPos] != 0x00 {
			scissor++
		}
		scissorPos--
	}

	return result[:scissorPos+1], nil
}
