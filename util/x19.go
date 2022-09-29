package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type X19User struct {
	Id    string
	Token string
}

type X19ResponseState struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var state X19ResponseState
	json.Unmarshal(data, &state)
	if state.Code != 0 {
		return nil, errors.New(state.Message)
	}

	return data, err
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
