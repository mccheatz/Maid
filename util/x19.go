package util

import (
	"bytes"
	"encoding/base64"
	"errors"
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

func X19ComputeDynamicToken(path string, body []byte, token string) string {
	var payload bytes.Buffer
	payload.WriteString(MD5Hex([]byte(token)))
	payload.Write(body)
	payload.WriteString("0eGsBkhl")
	payload.WriteString(path)

	sum := []byte(MD5Hex(payload.Bytes()))

	// convert the hex string to binary string by runes
	var binaryBuffer bytes.Buffer
	for _, by := range sum {
		a := 8
		b := 0x100
		for a != 0 {
			a--
			b = b >> 1
			if b&int(by) == 0 {
				binaryBuffer.WriteRune('0')
			} else {
				binaryBuffer.WriteRune('1')
			}
		}
	}

	// rotate the binary string
	r1 := binaryBuffer.String()
	binaryString := r1[6:] + r1[:6]

	// convert the binary string back and xor with the hex string
	for i := 0; i < len(sum); i++ {
		section := binaryString[i*8 : i*8+8]
		uVar14 := len(section)
		uVar9 := 0
		var by byte
		for uVar9 < len(section) {
			if section[uVar14-1] == '1' {
				by = by | 1<<(uVar9&0x1f)
			}
			uVar9++
			uVar14--
		}
		sum[i] = byte(by) ^ sum[i]
	}

	// encode the xor-ed hex string to base64 and only take first 16 bytes
	b64Encoded := base64.RawStdEncoding.EncodeToString(sum)
	resultReplacer := strings.NewReplacer("+", "m", "/", "o")
	result := resultReplacer.Replace(b64Encoded[:16] + "1")

	return result
}

func BuildX19Request(method string, address string, body []byte, userAgent string, user X19User) (*http.Request, error) {
	req, err := http.NewRequest(method, address, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", userAgent)

	// netease verify
	req.Header.Add("user-id", user.Id)
	{
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
		req.Header.Add("user-token", X19ComputeDynamicToken(path, body, user.Token))
	}

	return req, nil
}

func X19SimpleRequest(method string, url string, body []byte, client *http.Client, userAgent string, user X19User) ([]byte, error) {
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

func X19EncryptRequest(method string, address string, postBody []byte, client *http.Client, userAgent string, user X19User) ([]byte, error) {
	encryptedBody, err := X19HttpEncrypt(postBody)
	if err != nil {
		return nil, err
	}

	req, err := BuildX19Request(method, address, encryptedBody, userAgent, user)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	{
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
		req.Header.Set("user-token", X19ComputeDynamicToken(path, postBody, user.Token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return X19HttpDecrypt(body)
}
