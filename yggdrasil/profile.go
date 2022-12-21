package yggdrasil

import (
	"encoding/base64"
	"encoding/json"
)

type YggdrasilProfileResponse struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties []YggdrasilProfileProperty `json:"properties"`
}

type YggdrasilProfileProperty struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

type YggdrasilProfileTexture struct {
	Url string `json:"url"`
}

type YggdrasilProfileValue struct {
	Timestamp   uint64                             `json:"timestamp"`
	ProfileId   string                             `json:"profileId"`
	ProfileName string                             `json:"profileName"`
	Textures    map[string]YggdrasilProfileTexture `json:"textures"`
}

func (p *YggdrasilProfileValue) AddTexture(name, url string) {
	if p.Textures == nil {
		p.Textures = make(map[string]YggdrasilProfileTexture)
	}

	p.Textures[name] = YggdrasilProfileTexture{
		Url: url,
	}
}

func (p YggdrasilProfileValue) ToValue() (string, error) {
	jsonBody, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonBody), nil
}
