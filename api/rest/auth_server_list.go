package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type X19AuthServer struct {
	IP         string `json:"IP"`
	Port       int    `json:"Port"`
	ServerType string `json:"ServerType"`
}

func (server X19AuthServer) ToAddr() string {
	return server.IP + ":" + strconv.Itoa(server.Port)
}

func X19AuthServerList(client *http.Client, release X19ReleaseInfo, authServers *[]X19AuthServer) error {
	req, err := http.NewRequest("GET", release.AuthServerUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "WPFLauncher/0.0.0.0")

	resp1, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	body, err := io.ReadAll(resp1.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, authServers)
}
