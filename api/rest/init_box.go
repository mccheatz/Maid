package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

type BoxBasicInformation struct {
	GameId     string `json:"gid"`
	EncodeUrl  string `json:"url"`
	EncodeMode string `json:"mode"`
	T2S        bool   `json:"t2s"` // text to sound?
	F2H        bool   `json:"f2h"`
	Username   bool   `json:"un"`
	RStr       string `json:"rstr"` // replace str?
	Signature  bool   `json:"sig"`
	LU         bool   `json:"lu"`
	DRPF       string `json:"drpf"`
}

func InitBox(client *http.Client, info *BoxBasicInformation) error {
	req, err := http.NewRequest("GET", "http://optsdk.gameyw.netease.com/initbox_x19.html", nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "EnvSDK/1.0.9")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, info) // WHY DOES THIS RESPONSE A JSON LMAO
}
