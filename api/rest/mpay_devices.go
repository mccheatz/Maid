package rest

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"maid/util"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type MPayClientInfo struct {
	Brand         string `json:"brand"`
	DeviceModel   string `json:"device_model"`
	DeviceName    string `json:"device_name"`
	DeviceType    string `json:"device_type"`
	MacAddress    string `json:"mac"`
	Resolution    string `json:"resolution"`
	SystemName    string `json:"system_name"`
	SystemVersion string `json:"system_version"`
	Udid          string `json:"udid"`
	UniqueId      string `json:"unique_id"`
}

func (d *MPayClientInfo) GeneratePC() {
	d.Brand = "Microsoft"
	d.DeviceModel = "pc_mode"
	d.DeviceName = "DESKTOP-" + strings.ToUpper(util.RandStringRunes(7))
	d.DeviceType = "Computer"
	d.MacAddress = strings.ToUpper(util.RandMacAddress())
	d.Resolution = "1920*1080"
	d.SystemName = "windows"
	d.SystemVersion = "10"
	d.Udid = strings.ReplaceAll(uuid.NewString(), "-", "")
	d.UniqueId = strings.ReplaceAll(uuid.NewString(), "-", "")
}

func (d *MPayClientInfo) GenerateMobile() {
	brands := []string{"Xiaomi", "Huawei", "HTC", "Google", "Nokia", "Sharp", "Asus"}
	d.Brand = brands[rand.Intn(len(brands)-1)]
	d.DeviceName = strings.ToUpper(util.RandStringRunes(3)) + "-AL00" // like huawei phone models cuz this easy to generate
	d.DeviceModel = d.DeviceName
	d.DeviceType = "mobile"
	d.MacAddress = strings.ToUpper(util.RandMacAddress())
	d.Resolution = "1080*2027" // 2k screen
	d.SystemName = "Android"
	d.SystemVersion = strconv.Itoa(7 + rand.Intn(5))
	d.Udid = strings.ReplaceAll(uuid.NewString(), "-", "")
	d.UniqueId = strings.ReplaceAll(uuid.NewString(), "-", "")
}

type MPayAppInfo struct {
	AppMode             string `json:"app_mode"`
	AppType             string `json:"app_type"`
	Arch                string `json:"arch"`
	ClientVersion       string `json:"cv"`
	GameId              string `json:"game_id"`
	GameVersion         string `json:"gv"`
	MCountAppKey        string `json:"mcount_app_key"`
	MCountTransactionId string `json:"mcount_transaction_id"`
	OptFields           string `json:"opt_fields"`
	ProcessId           string `json:"process_id"`
	ServiceVersion      string `json:"sv"`
	UpdaterVersion      string `json:"updater_cv"`
	GVN                 string `json:"gvn"`
	CloudExtraBase64    string `json:"_cloud_extra_base64"`
}

func (mp *MPayAppInfo) GenerateForX19(version string) {
	mp.AppMode = "2"
	mp.AppType = "games"
	mp.Arch = "win_x32"
	mp.ClientVersion = "c3.4.0"
	mp.GameId = "aecfrxodyqaaaajp-g-x19"
	mp.GameVersion = version
	mp.MCountAppKey = "EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq"
	mp.MCountTransactionId = uuid.NewString() + "-2"
	mp.OptFields = "nickname,avatar,realname_status,mobile_bind_status"
	mp.ProcessId = strconv.Itoa(1000 + rand.Intn(10000))
	mp.ServiceVersion = "10"
	mp.UpdaterVersion = "c1.0.0"
}

func (mp *MPayAppInfo) GenerateForX19Mobile(version string) {
	mp.AppMode = "2"
	mp.AppType = "games"
	mp.ClientVersion = "a3.32.1"
	mp.GameId = "aecfrxodyqaaaajp-g-x19"
	mp.GameVersion = version
	mp.MCountAppKey = "EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq"
	mp.MCountTransactionId = uuid.NewString() + "-2"
	mp.OptFields = "nickname,avatar,realname_status,mobile_bind_status,exit_popup_info"
	mp.ServiceVersion = "31"
	mp.GVN = "2.2.15.204111"
	mp.CloudExtraBase64 = "eyJleHRyYSI6e319" // {"extra":{}}
}

type MPayError struct {
	Reason string `json:"reason"`
	Code   int    `json:"code"`
}

func (e MPayError) Error() string {
	return e.Reason + "(code=" + strconv.Itoa(e.Code) + ")"
}

func mPayErrorHandle(statusCode int, body []byte) error {
	if statusCode/200 == 2 {
		var errReason MPayError
		err := json.Unmarshal(body, &errReason)
		if err != nil {
			errReason.Code = statusCode
			errReason.Reason = err.Error()
		}
		return errReason
	}

	return nil
}

type MPayDevice struct {
	Id        string `json:"id"`
	Key       string `json:"key"`
	BinaryKey []byte
}

func (md *MPayDevice) ClaimBinaryKey() error {
	key, err := hex.DecodeString(md.Key)
	if err == nil {
		md.BinaryKey = key
	}
	return err
}

func MPayDevices(client *http.Client, clientMPay MPayClientInfo, appMPay MPayAppInfo, device *MPayDevice) error {
	postBody := url.Values{}

	util.PushToParameters(clientMPay, &postBody)
	util.PushToParameters(appMPay, &postBody)

	req, err := http.NewRequest("POST", "https://service.mkey.163.com/mpay/"+appMPay.AppType+"/"+appMPay.GameId+"/devices", bytes.NewBuffer([]byte(postBody.Encode())))
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

	var query map[string]MPayDevice

	err = json.Unmarshal(body, &query)
	if err != nil {
		return err
	}

	if val, ok := query["device"]; ok {
		*device = val
		return nil
	}

	return errors.New("no device info found in response")
}
