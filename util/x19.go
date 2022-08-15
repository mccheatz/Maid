package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type X19User struct {
	Id    string
	Token string
}

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

type test1 struct {
	a int
	b int
	f [0xf]byte
}

func (a *test1) c(b byte) {
	a.a = 0
	a.b = 0xf

	uVar3 := 8
	uVar4 := 0x100
	for uVar3 != 0 {
		uVar3 = uVar3 - 1
		uVar4 = uVar4>>1 | BoolToInt((uVar4&1) != 0)<<0x1f
		if b&byte(uVar4) == 0 {
			a.d(0x30)
		} else {
			a.d(0x31)
		}
	}
}

func (a *test1) d(b byte) {
	c := a.a
	d := a.b
	if c < d {
		a.a = c + 1
		a.f[c] = b
		// println(a.a)
		// println(a.b - a.a)
		a.f[c+1] = 0
		return
	}

	panic(errors.New("=w="))
}

func X19SpecialMD5(body []byte) []byte {
	sum := md5.Sum(body)
	for i := 0; i < len(sum); i++ {
		b := sum[i]
		if b >= 'A' && b <= 'Z' {
			sum[i] += 'a' - 'A'
		}
	}

	return sum[:]
}

func ComputeDynamicToken(path string, body []byte, token string) string {
	tokenMd5 := X19SpecialMD5([]byte(token))
	tokenMd5 = append(tokenMd5, body...)
	tokenMd5 = append(tokenMd5, []byte("0eGsBkhl")...)
	tokenMd5 = append(tokenMd5, []byte(path)...)

	mergedMd5 := X19SpecialMD5(tokenMd5)

	b, _ := base64.StdEncoding.DecodeString("o7m7mu49ro9prqor1")

	for i := 0; i < len(mergedMd5); i++ {
		if i < len(b) {
			fmt.Printf("%x", mergedMd5[i]^b[i])
		} else {
			print("-")
		}
		print("\t")
		fmt.Printf("%x\n", mergedMd5[i])
	}

	return ""

	local_b8 := make([]byte, 0)

	for i := 0; i < len(mergedMd5); i++ {
		a := test1{}
		a.c(mergedMd5[i])
		for j := 0; j < a.a; j++ {
			local_b8 = append(local_b8, a.f[j])
		}
	}

	// for i := 0; i < len(local_b8); i++ {
	// 	println(local_b8[i])
	// }

	processedPayload := make([]byte, len(local_b8))
	for i := 0; i < len(local_b8)-6; i++ {
		processedPayload[i] = local_b8[6+i]
	}
	for i := 0; i < 6; i++ {
		processedPayload[len(local_b8)-6+i] = local_b8[i]
	}
	for i := 0; i < len(mergedMd5); i++ {
		processedPayload[i] = processedPayload[i] ^ mergedMd5[i]
	}

	println(hex.EncodeToString(mergedMd5))
	println(hex.EncodeToString(processedPayload))

	b64Encoded := base64.RawStdEncoding.EncodeToString(processedPayload)
	resultReplacer := strings.NewReplacer("+", "m", "/", "o")
	result := resultReplacer.Replace(b64Encoded[:16] + "1")

	return result
}

func BuildX19Request(method string, address string, body []byte, userAgent string, user *X19User) (*http.Request, error) {
	req, err := http.NewRequest(method, address, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)
	if user != nil {
		req.Header.Add("user-id", user.Id)
		// TODO user-token
		u, err := url.Parse(address)
		if err != nil {
			panic(err)
		}
		path := u.Path
		if len(u.RawQuery) != 0 {
			path += "?" + u.RawQuery
		}
		if len(u.Fragment) != 0 {
			path += "#" + u.Fragment
		}
		req.Header.Add("user-token", ComputeDynamicToken(path, body, user.Token))
	}

	return req, nil
}

func X19SimpleRequest(method string, url string, body []byte, client *http.Client, userAgent string, user *X19User) ([]byte, error) {
	req, err := BuildX19Request(method, url, body, userAgent, user)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
