package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

type X19ReleaseInfo struct {
	HostNum                    int
	ServerHostNum              int
	TempServerStop             int
	ServerStop                 string
	CdnUrl                     string
	StaticWebVersionUrl        string
	SeadraUrl                  string
	EmbedWebPageUrl            string
	NewsVideo                  string
	GameCenter                 string
	VideoPrefix                string
	ComponentCenter            string
	GameDetail                 string
	CompDetail                 string
	LiveUrl                    string
	ForumUrl                   string
	WebServerUrl               string
	WebServerGrayUrl           string
	CoreServerUrl              string
	TransferServerUrl          string
	PeTransferServerUrl        string
	PeTransferServerHttpUrl    string
	TransferServerHttpUrl      string
	PeTransferServerNewHttpUrl string
	AuthServerUrl              string
	AuthServerCppUrl           string
	AuthorityUrl               string
	CustomerServiceUrl         string
	ChatServerUrl              string
	PathNUrl                   string
	PePathNUrl                 string
	MgbSdkUrl                  string
	DCWebUrl                   string
	ApiGatewayUrl              string
	ApiGatewayGrayUrl          string
	PlatformUrl                string
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
