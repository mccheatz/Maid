package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

type NeteaseAuthServer struct {
	IP         string
	Post       int
	ServerType string
}

func AuthServerList(client *http.Client, release NeteaseReleaseInfo, authServers *[]NeteaseAuthServer) error {
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
