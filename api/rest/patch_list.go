package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type X19Patch struct {
	Name string
	Size int    `json:"size"`
	Url  string `json:"url"`
	Md5  string `json:"md5"`
}

type X19Version struct {
	Version     string `json:"version"`
	LauncherMD5 string `json:"launcher_md5"`
	UpdaterMD5  string `json:"updater_md5"`
}

func X19PatchList(client *http.Client, patchList *[]X19Patch) error {
	req, err := http.NewRequest("GET", "https://x19.update.netease.com/pl/x19_java_patchlist", nil)
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

	body = []byte("{" + strings.TrimSuffix(strings.TrimSpace(string(body)), ",") + "}")

	var patches map[string]X19Patch

	err = json.Unmarshal(body, &patches)
	if err != nil {
		return err
	}

	// merge key into value
	*patchList = make([]X19Patch, 0)
	for version, patch := range patches {
		patch.Name = version

		*patchList = append(*patchList, patch)
	}

	return nil
}
