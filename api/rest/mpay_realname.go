package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"maid/util"
	"net/http"
	"net/url"
)

type MPayRealNameResult struct {
	NeedAAS      bool   `json:"need_aas"`
	RealNameType string `json:"realname_type"`
}

func mPayRealNameBase(client *http.Client, device MPayDevice, appMPay MPayAppInfo, user MPayUser, realname string, idRegion string, idNum string, result *MPayRealNameResult, subApi string) error {
	postBody := url.Values{}

	util.PushToParameters(appMPay, &postBody)
	postBody.Add("device_id", device.Id)
	postBody.Add("user_id", user.Id)
	postBody.Add("token", user.Token)
	postBody.Add("realname", realname)
	postBody.Add("id_region", idRegion)
	postBody.Add("id_num", idNum)

	req, err := http.NewRequest("POST", "https://service.mkey.163.com/mpay/api/users/realname/"+subApi, bytes.NewBuffer([]byte(postBody.Encode())))
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

	if result != nil {
		return json.Unmarshal(body, &result)
	}

	return nil
}

// This only upload the realname information to the server and verify is the information able to use or not.
func MPayRealNameVerify(client *http.Client, device MPayDevice, appMPay MPayAppInfo, user MPayUser, realname string, idRegion string, idNum string, result *MPayRealNameResult) error {
	return mPayRealNameBase(client, device, appMPay, user, realname, idRegion, idNum, result, "verify")
}

func MPayRealNameUpdate(client *http.Client, device MPayDevice, appMPay MPayAppInfo, user MPayUser, realname string, idRegion string, idNum string, result *MPayRealNameResult) error {
	return mPayRealNameBase(client, device, appMPay, user, realname, idRegion, idNum, result, "update_by_token")
}
