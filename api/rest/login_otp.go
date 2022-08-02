package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type X19OTPEntity struct {
	Code    int    `json:"code"`
	Details string `json:"details"`
	Entity  struct {
		AId      int    `json:"aid"`
		LockTime int    `json:"lock_time"`
		OpenOTP  int    `json:"open_otp"`
		OTP      int    `json:"otp"`
		OTPToken string `json:"otp_token"`
	} `json:"entity"`
	Message string `json:"message"`
}

func LoginOTP(client *http.Client, sAuth MPaySAuthToken, userAgent string, coreServerUrl string, otpEntity *X19OTPEntity) error {
	sAuthJson, err := json.Marshal(sAuth)
	if err != nil {
		return err
	}
	sAuthContainer := struct {
		SAuthJson string `json:"sauth_json"`
	}{
		SAuthJson: string(sAuthJson),
	}

	sAuthJson, err = json.Marshal(sAuthContainer)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", coreServerUrl+"/login-otp", bytes.NewBuffer(sAuthJson))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &otpEntity)
}
