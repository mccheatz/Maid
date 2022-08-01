package rest

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maid/util"
	"net/http"
	"net/url"
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
	RealnameStatus       int    `json:"realname_status"`
	RealnameVerifyStatus int    `json:"realname_verify_status"`
	Token                string `json:"token"`
}

func mPayLoginParams(username string, password string, device MPayDevice, client MPayClientInfo) (string, error) {
	err := device.ClaimBinaryKey()
	if err != nil {
		return "", err
	}

	unencrypted, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
		UniqueId string `json:"unique_id"`
	}{
		Username: username,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(password))),
		UniqueId: client.UniqueId, // is this right?
	})
	if err != nil {
		return "", err
	}

	encrypted, err := util.AesPkcs7Encrypt(device.BinaryKey, unencrypted)
	if err != nil {
		return "", err
	}

	return util.ToHexString(encrypted), nil
}

func MPayLogin(client *http.Client, device MPayDevice, appMPay MPayAppInfo, clientMPay MPayClientInfo, username string, password string, user *MPayUser) error {
	postBody := url.Values{}

	params, err := mPayLoginParams(username, password, device, clientMPay)
	if err != nil {
		return err
	}

	util.PushToParameters(appMPay, &postBody)
	postBody.Add("un", base64.StdEncoding.EncodeToString([]byte(username)))
	postBody.Add("params", params)
	postBody.Add("app_channel", "netease")

	req, err := http.NewRequest("POST", "https://service.mkey.163.com/mpay/games/aecfrxodyqaaaajp-g-x19/devices/"+device.Id+"/users", bytes.NewBuffer([]byte(postBody.Encode())))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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
