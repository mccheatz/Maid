package rest

import (
	"encoding/json"
	"maid/util"
	"net/http"
)

type LauncherPersonInfo struct {
	Code    int                      `json:"code"`
	Details string                   `json:"details"`
	Entity  LauncherPersonInfoEntity `json:"entity"`
	Message string                   `json:"message"`
}

type LauncherPersonInfoEntity struct {
	EntityId string `json:"entityId"`
	Gender   string `json:"gender"`
	// avatar
	HeadImage string `json:"headImage"`
	Nickname  string `json:"nickname"`
	// bio
	Signature string `json:"signature"`
}

func FetchLauncherPersonInfo(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, result *LauncherPersonInfo) error {
	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/personal-info/get", []byte{}, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

type LauncherSetNicknameRequest struct {
	Name     string `json:"name"`
	EntityId string `json:"entityId"`
}

type LauncherSetNicknameResponse struct {
	Code    int                        `json:"code"`
	Details string                     `json:"details"`
	Entity  LauncherSetNicknameRequest `json:"entity"`
	Message string                     `json:"message"`
}

func LauncherSetNickname(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, request LauncherSetNicknameRequest, result *LauncherSetNicknameResponse) error {
	postBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/nickname-setting", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}
