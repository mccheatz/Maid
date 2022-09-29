package mcproto

import (
	"encoding/json"
	"maid/util"

	"github.com/Tnze/go-mc/bot"
)

type PingResponse struct {
	Description util.JsonRaw `json:"descrption"`
	Players     struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	} `json:"players"`
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"version"`
	} `json:"version"`
	ModInfo struct {
		Type string `json:"type"`
	} `json:"modinfo"`
}

func PingServer(addr string, response *PingResponse) error {
	resp, _, err := bot.PingAndList(addr)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return err
	}

	return nil
}
