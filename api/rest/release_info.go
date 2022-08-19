package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

type X19ReleaseInfo struct {
	HostNum                    int    `json:"HostNum"`
	ServerHostNum              int    `json:"ServerHostNum"`
	TempServerStop             int    `json:"TempServerStop"`
	ServerStop                 string `json:"ServerStop"`
	CdnUrl                     string `json:"CdnUrl"`
	StaticWebVersionUrl        string `json:"StaticWebVersionUrl"`
	SeadraUrl                  string `json:"SeadraUrl"`
	EmbedWebPageUrl            string `json:"EmbedWebPageUrl"`
	NewsVideo                  string `json:"NewsVideo"`
	GameCenter                 string `json:"GameCenter"`
	VideoPrefix                string `json:"VideoPrefix"`
	ComponentCenter            string `json:"ComponentCenter"`
	GameDetail                 string `json:"GameDetail"`
	CompDetail                 string `json:"CompDetail"`
	LiveUrl                    string `json:"LiveUrl"`
	ForumUrl                   string `json:"ForumUrl"`
	WebServerUrl               string `json:"WebServerUrl"`
	WebServerGrayUrl           string `json:"WebServerGrayUrl"`
	CoreServerUrl              string `json:"CoreServerUrl"`
	TransferServerUrl          string `json:"TransferServerUrl"`
	PeTransferServerUrl        string `json:"PeTransferServerUrl"`
	PeTransferServerHttpUrl    string `json:"PeTransferServerHttpUrl"`
	TransferServerHttpUrl      string `json:"TransferServerHttpUrl"`
	PeTransferServerNewHttpUrl string `json:"PeTransferServerNewHttpUrl"`
	AuthServerUrl              string `json:"AuthServerUrl"`
	AuthServerCppUrl           string `json:"AuthServerCppUrl"`
	AuthorityUrl               string `json:"AuthorityUrl"`
	CustomerServiceUrl         string `json:"CustomerServiceUrl"`
	ChatServerUrl              string `json:"ChatServerUrl"`
	PathNUrl                   string `json:"PathNUrl"`
	PePathNUrl                 string `json:"PePathNUrl"`
	MgbSdkUrl                  string `json:"MgbSdkUrl"`
	DCWebUrl                   string `json:"DCWebUrl"`
	ApiGatewayUrl              string `json:"ApiGatewayUrl"`
	ApiGatewayGrayUrl          string `json:"ApiGatewayGrayUrl"`
	PlatformUrl                string `json:"PlatformUrl"`
}

func X19ReleaseInfoFetch(client *http.Client, release *X19ReleaseInfo) error {
	req, err := http.NewRequest("GET", "https://x19.update.netease.com/serverlist/release.json", nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "WPFLauncher/0.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, release)
}
