package rest

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"maid/util"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type MPayUser struct {
	Avatar            string `json:"avatar"`
	ClientUsername    string `json:"client_username"`
	DisplayerUsername string `json:"display_username"`
	ExtAccessToken    string `json:"ext_access_token"`
	Id                string `json:"id"`
	LoginChannel      string `json:"login_channel"`
	LoginType         int    `json:"login_type"`
	MobileBindStatus  int    `json:"mobile_bind_status"`
	NeedAAS           bool   `json:"need_aas"`
	NeedMask          bool   `json:"need_mask"`
	Nickname          string `json:"nickname"`
	PcExtInfo         struct {
		ExtraUnisdkData string `json:"ext_unisdk_data"`
		FromGameId      string `json:"from_game_id"`
		AppChannel      string `json:"src_app_channel"`
		ClientIp        string `json:"src_client_ip"`
		ClientType      int    `json:"src_client_type"`
		JfGameId        string `json:"src_jf_game_id"`
		PayChannel      string `json:"src_pay_channel"`
		SdkVersion      string `json:"src_sdk_version"`
		Udid            string `json:"src_udid"`
	} `json:"pc_ext_info"`
	RealNameStatus       int    `json:"realname_status"`
	RealNameVerifyStatus int    `json:"realname_verify_status"`
	Token                string `json:"token"`
}

type MPaySAuthToken struct {
	GameId         string `json:"gameid"`
	LoginChannel   string `json:"login_channel"`
	AppChannel     string `json:"app_channel"`
	Platform       string `json:"platform"`
	SdkUid         string `json:"sdkuid"`
	SessionId      string `json:"sessionid"`
	SdkVersion     string `json:"sdk_version"`
	Udid           string `json:"udid"`
	DeviceId       string `json:"deviceid"`
	AimInfo        string `json:"aim_info"`
	ClientLoginSn  string `json:"client_login_sn"`
	GasToken       string `json:"gas_token"`
	SourcePlatform string `json:"source_platform"`
	Ip             string `json:"ip"`
}

func (mp *MPayUser) ConvertToSAuth(id string, client MPayClientInfo, device MPayDevice) MPaySAuthToken {
	sa := MPaySAuthToken{}
	sa.GameId = id
	sa.LoginChannel = "netease"
	sa.AppChannel = "netease"
	sa.Platform = "pc"
	sa.SdkUid = mp.Id
	sa.SessionId = mp.Token
	sa.SdkVersion = "3.4.0"
	sa.Udid = client.Udid
	sa.DeviceId = device.Id
	sa.AimInfo = "{\"aim\":\"" + mp.PcExtInfo.ClientIp + "\",\"country\":\"CN\",\"tz\":\"+0800\",\"tzid\":\"\"}"
	sa.ClientLoginSn = strings.ToUpper(strings.ReplaceAll(uuid.NewString(), "-", ""))
	sa.GasToken = ""
	sa.SourcePlatform = "pc"
	sa.Ip = mp.PcExtInfo.ClientIp

	return sa
}

func mPayEncryptToParams(unencrypted []byte, device MPayDevice) (string, error) {
	err := device.ClaimBinaryKey()
	if err != nil {
		return "", err
	}

	encrypted, err := util.AES_ECB_PKCS7Encrypt(device.BinaryKey, unencrypted)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}

func mPayLoginBase(client *http.Client, appMPay MPayAppInfo, params string, postUrl string, user *MPayUser) error {
	postBody := url.Values{}

	util.PushToParameters(appMPay, &postBody)
	postBody.Add("params", params)
	postBody.Add("app_channel", "netease")

	req, err := http.NewRequest("POST", "https://service.mkey.163.com/mpay/"+postUrl, bytes.NewBuffer([]byte(postBody.Encode())))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = mPayErrorHandle(resp.StatusCode, body)
	if err != nil {
		return err
	}

	var query map[string]MPayUser

	err = json.Unmarshal(body, &query)
	if err != nil {
		return err
	}

	if val, ok := query["user"]; ok {
		*user = val
		return nil
	}

	return errors.New("no device info found in response")
}

func MPayLogin(client *http.Client, device MPayDevice, appMPay MPayAppInfo, clientMPay MPayClientInfo, username string, password string, user *MPayUser) error {
	unencrypted, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
		UniqueId string `json:"unique_id"`
	}{
		Username: username,
		Password: hex.EncodeToString(util.MD5Sum([]byte(password))),
		UniqueId: clientMPay.UniqueId,
	})
	if err != nil {
		return err
	}

	params, err := mPayEncryptToParams(unencrypted, device)
	if err != nil {
		return err
	}

	return mPayLoginBase(client, appMPay, params, appMPay.AppType+"/"+appMPay.GameId+"/devices/"+device.Id+"/users", user)
}

// this will only works in games that support guest login
func MPayLoginGuest(client *http.Client, device MPayDevice, appMPay MPayAppInfo, clientMPay MPayClientInfo, user *MPayUser) error {
	unencrypted, err := json.Marshal(struct {
		Udid string `json:"udid"`
	}{
		Udid: clientMPay.Udid,
	})
	if err != nil {
		return err
	}

	params, err := mPayEncryptToParams(unencrypted, device)
	if err != nil {
		return err
	}

	return mPayLoginBase(client, appMPay, params, appMPay.AppType+"/"+appMPay.GameId+"/devices/"+device.Id+"/users/by_guest", user)
}
