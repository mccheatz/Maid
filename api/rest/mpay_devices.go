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
	InitUrsDevice string `json:"init_urs_device"`
	MacAddress    string `json:"mac"`
	Resolution    string `json:"resolution"`
	SystemName    string `json:"system_name"`
	SystemVersion string `json:"system_version"`
	Udid          string `json:"udid"`
	UniqueId      string `json:"unique_id"`
}

func (d *MPayClientInfo) Generate() {
	d.Brand = "Microsoft"
	d.DeviceModel = "pc_mode"
	d.DeviceName = "DESKTOP-" + strings.ToUpper(util.RandStringRunes(7))
	d.DeviceType = "Computer"
	d.InitUrsDevice = "0"
	d.MacAddress = strings.ToUpper(util.RandMacAddress())
	d.Resolution = "1920*1080"
	d.SystemName = "windows"
	d.SystemVersion = "10"
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

	req, err := http.NewRequest("POST", "https://service.mkey.163.com/mpay/games/aecfrxodyqaaaajp-g-x19/devices", bytes.NewBuffer([]byte(postBody.Encode())))
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
