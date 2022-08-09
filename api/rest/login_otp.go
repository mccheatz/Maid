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

func LoginOTP(client *http.Client, sAuth MPaySAuthToken, userAgent string, release X19ReleaseInfo, otpEntity *X19OTPEntity) error {
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

func GenerateAuthenticationOTPBody(client *http.Client, userAgent string, sAuth MPaySAuthToken, clientMpay MPayClientInfo, release X19ReleaseInfo, version X19Version, otpEntity X19OTPEntity) error {
	// sAuthJson, err := json.Marshal(sAuth)
	// if err != nil {
	// 	return err
	// }
	// saData := struct {
	// 	OsName       string `json:"os_name"`
	// 	OsVersion    string `json:"os_ver"`
	// 	MacAddress   string `json:"mac_addr"`
	// 	Udid         string `json:"udid"`
	// 	AppVersion   string `json:"app_ver"`
	// 	SdkVersion   string `json:"sdk_ver"`
	// 	Network      string `json:"network"`
	// 	Disk         string `json:"disk"`
	// 	Is64Bit      string `json:"is64bit"`
	// 	VideoCard1   string `json:"video_card1"`
	// 	VideoCard2   string `json:"video_card2"`
	// 	VideoCard3   string `json:"video_card3"`
	// 	VideoCard4   string `json:"video_card4"`
	// 	LauncherType string `json:"launcher_type"`
	// 	PayChannel   string `json:"pay_channel"`
	// }{
	// 	OsName:       "windows",
	// 	OsVersion:    "Microsoft Windows 10",
	// 	MacAddress:   clientMpay.MacAddress,
	// 	Udid:         clientMpay.Udid,
	// 	AppVersion:   "0.0.0.0",
	// 	Disk:         fmt.Sprintf("%02x%02x%02x%02x", rand.Intn(0xff), rand.Intn(0xff), rand.Intn(0xff), rand.Intn(0xff)),
	// 	Is64Bit:      "1",
	// 	VideoCard1:   "Nvidia GTX 1080 Ti",
	// 	LauncherType: "PC_java",
	// 	PayChannel:   "netease",
	// }
	// saDataJson, err := json.Marshal(saData)
	// if err != nil {
	// 	return err
	// }

	// bodyStruct := struct {
	// 	SaData           string        `json:"sa_data"`
	// 	SAuthJson        string        `json:"sauth_json"`
	// 	Version          X19Version    `json:"version"`
	// 	SdkUid           util.JsonNull `json:"sdkuid"`
	// 	AId              string        `json:"aid"`
	// 	HasMessage       bool          `json:"hasMessage"` // imagine tracker in game
	// 	HasGMail         bool          `json:"hasGmail"`
	// 	OTPToken         string        `json:"otp_token"`
	// 	OTPPassword      util.JsonNull `json:"otp_pwd"`
	// 	LockTime         int           `json:"lock_time"`
	// 	Env              util.JsonNull `json:"env"`
	// 	MinEngineVersion util.JsonNull `json:"min_engine_version"`
	// 	MinPatchVersion  util.JsonNull `json:"min_patch_version"`
	// 	VerifyStatus     int           `json:"verify_status"`
	// 	UniSDKLoginJson  util.JsonNull `json:"unisdk_login_json"`
	// 	EntityId         util.JsonNull `json:"entity_id"`
	// }{
	// 	SaData:       string(saDataJson),
	// 	SAuthJson:    string(sAuthJson),
	// 	Version:      version,
	// 	AId:          strconv.Itoa(otpEntity.Entity.AId),
	// 	HasMessage:   false,
	// 	HasGMail:     false,
	// 	OTPToken:     otpEntity.Entity.OTPToken,
	// 	LockTime:     0,
	// 	VerifyStatus: 0,
	// }

	// postBody, err := json.Marshal(bodyStruct)
	// if err != nil {
	// 	return err
	// }

	// postBody := []byte(`{"sa_data":"{\"os_name\":\"windows\",\"os_ver\":\"Microsoft Windows 10\",\"mac_addr\":\"8641E01A2A09\",\"udid\":\"457759472c844f80ae8080fbb824f042\",\"app_ver\":\"0.0.0.0\",\"sdk_ver\":\"\",\"network\":\"\",\"disk\":\"c29a6378\",\"is64bit\":\"1\",\"video_card1\":\"Nvidia GTX 1080 Ti\",\"video_card2\":\"\",\"video_card3\":\"\",\"video_card4\":\"\",\"launcher_type\":\"PC_java\",\"pay_channel\":\"netease\"}","sauth_json":"{\"gameid\":\"x19\",\"login_channel\":\"netease\",\"app_channel\":\"netease\",\"platform\":\"pc\",\"sdkuid\":\"aebgdosotvr757nu\",\"sessionid\":\"1-eyJzIjogIjg1ZTYyOTM2YmIxODRiYjQ0NjVmODhiZDQ4ZWM4MTdkODY5OGZmIiwgImdfaSI6ICJhZWNmcnhvZHlxYWFhYWpwIiwgInQiOiAxfSAg\",\"sdk_version\":\"3.4.0\",\"udid\":\"457759472c844f80ae8080fbb824f042\",\"deviceid\":\"amawf4iaakwwwp5q-d\",\"aim_info\":\"{\\\"aim\\\":\\\"\\\",\\\"country\\\":\\\"CN\\\",\\\"tz\\\":\\\"+0800\\\",\\\"tzid\\\":\\\"\\\"}\",\"client_login_sn\":\"A314DF9A462545C58578B403736ABAA3\",\"gas_token\":\"\",\"source_platform\":\"pc\",\"ip\":\"\"}","version":{"version":"1.8.21.53078","launcher_md5":"","updater_md5":""},"sdkuid":"","aid":"595682882","hasMessage":false,"hasGmail":false,"otp_token":"CoIh3LcAMQZOdu4L","otp_pwd":"","lock_time":0,"env":"","min_engine_version":"","min_patch_version":"","verify_status":0,"unisdk_login_json":"","entity_id":""}`)

	// a := struct {
	// 	Body string `json:"body"`
	// }{
	// 	Body: string(postBody),
	// }

	// b, err := json.Marshal(a)
	// if err != nil {
	// 	return err
	// }

	// println(string(b))

	// postBody, err := util.X19HttpEncrypt(postBody)
	// if err != nil {
	// 	return err
	// }

	// println(base64.StdEncoding.EncodeToString(postBody))

	postBody := []byte{}

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

	println(len(body))

	return nil
}
