package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maid/util"
	"math/rand"
	"net/http"
	"strconv"
)

type X19OTPEntity struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Entity  struct {
		AId      int    `json:"aid"`
		LockTime int    `json:"lock_time"`
		OpenOTP  int    `json:"open_otp"`
		OTP      int    `json:"otp"`
		OTPToken string `json:"otp_token"`
	} `json:"entity"`
}

func X19LoginOTP(client *http.Client, sAuth MPaySAuthToken, userAgent string, release X19ReleaseInfo, otpEntity *X19OTPEntity) error {
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

	req, err := http.NewRequest("POST", release.CoreServerUrl+"/login-otp", bytes.NewBuffer(sAuthJson))
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

type X19AuthenticationEntity struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Entity  struct {
		EntityId        string `json:"entity_id"`
		Account         string `json:"account"`
		Token           string `json:"token"`
		Sead            string `json:"sead"`
		HasMessage      bool   `json:"hasMessage"`
		HasGmail        bool   `json:"hasGmail"`
		AId             string `json:"aid"`
		SdkUid          string `json:"sdkuid"`
		AccessToken     string `json:"access_token"`
		UniSDKLoginJson string `json:"unisdk_login_json"`
		VerifyStatus    int    `json:"verify_status"`
		IsRegister      bool   `json:"is_register"`
		// "autopatch": []
		Env              string `json:"env"`
		LastServerUpTime int    `json:"last_server_up_time"`
		MinEngineVersion string `json:"min_engine_version"`
		MinPatchVersion  string `json:"min_patch_version"`
	} `json:"entity"`
}

func X19AuthenticationOTP(client *http.Client, userAgent string, sAuth MPaySAuthToken, clientMpay MPayClientInfo, release X19ReleaseInfo, version X19Version, otpEntity X19OTPEntity, authEntity *X19AuthenticationEntity) error {
	sAuthJson, err := json.Marshal(sAuth)
	if err != nil {
		return err
	}
	saData := struct {
		OsName       string `json:"os_name"`
		OsVersion    string `json:"os_ver"`
		MacAddress   string `json:"mac_addr"`
		Udid         string `json:"udid"`
		AppVersion   string `json:"app_ver"`
		SdkVersion   string `json:"sdk_ver"`
		Network      string `json:"network"`
		Disk         string `json:"disk"`
		Is64Bit      string `json:"is64bit"`
		VideoCard1   string `json:"video_card1"`
		VideoCard2   string `json:"video_card2"`
		VideoCard3   string `json:"video_card3"`
		VideoCard4   string `json:"video_card4"`
		LauncherType string `json:"launcher_type"`
		PayChannel   string `json:"pay_channel"`
	}{
		OsName:       "windows",
		OsVersion:    "Microsoft Windows 10",
		MacAddress:   clientMpay.MacAddress,
		Udid:         clientMpay.Udid,
		AppVersion:   "0.0.0.0",
		Disk:         fmt.Sprintf("%02x%02x%02x%02x", rand.Intn(0xff), rand.Intn(0xff), rand.Intn(0xff), rand.Intn(0xff)),
		Is64Bit:      "1",
		VideoCard1:   "Nvidia GTX 1080 Ti",
		LauncherType: "PC_java",
		PayChannel:   "netease",
	}
	saDataJson, err := json.Marshal(saData)
	if err != nil {
		return err
	}

	bodyStruct := struct {
		SaData           string        `json:"sa_data"`
		SAuthJson        string        `json:"sauth_json"`
		Version          X19Version    `json:"version"`
		SdkUid           util.JsonNull `json:"sdkuid"`
		AId              string        `json:"aid"`
		HasMessage       bool          `json:"hasMessage"` // imagine tracker in game
		HasGMail         bool          `json:"hasGmail"`
		OTPToken         string        `json:"otp_token"`
		OTPPassword      util.JsonNull `json:"otp_pwd"`
		LockTime         int           `json:"lock_time"`
		Env              util.JsonNull `json:"env"`
		MinEngineVersion util.JsonNull `json:"min_engine_version"`
		MinPatchVersion  util.JsonNull `json:"min_patch_version"`
		VerifyStatus     int           `json:"verify_status"`
		UniSDKLoginJson  util.JsonNull `json:"unisdk_login_json"`
		EntityId         util.JsonNull `json:"entity_id"`
	}{
		SaData:       string(saDataJson),
		SAuthJson:    string(sAuthJson),
		Version:      version,
		AId:          strconv.Itoa(otpEntity.Entity.AId),
		HasMessage:   false,
		HasGMail:     false,
		OTPToken:     otpEntity.Entity.OTPToken,
		LockTime:     0,
		VerifyStatus: 0,
	}

	postBody, err := json.Marshal(bodyStruct)
	if err != nil {
		return err
	}

	postBody, err = util.X19HttpEncrypt(postBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", release.CoreServerUrl+"/authentication-otp", bytes.NewBuffer(postBody))
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

	body, err = util.X19HttpDecrypt(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &authEntity)
}
