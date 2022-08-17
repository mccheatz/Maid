package rest

import (
	"encoding/json"
	"maid/util"
	"net/http"
)

const (
	CHARACTER_DELETE_MODE_PRE    string = "pre-delete"
	CHARACTER_DELETE_MODE_CANCEL string = "cancel-pre-delete"
	CHARACTER_DELETE_MODE_DELETE string = "delete"
)

type X19GameCharacterQueryInfo struct {
	GameId   string `json:"game_id"`
	GameType int    `json:"game_type"`
	Length   int    `json:"length"`
	Offset   int    `json:"offset"`
	UserId   string `json:"user_id"`
}

type X19CreateGameCharacterInfo struct {
	GameId   string `json:"game_id"`
	GameType int    `json:"game_type"`
	Name     string `json:"name"`
	UserId   string `json:"user_id"`
}

type X19GameCharacterQueryEntity struct {
	EntityId   string `json:"entity_id"`
	CreateTime int64  `json:"create_time"`
	ExpireTime int64  `json:"expire_time"`
	GameId     string `json:"game_id"`
	GameType   int    `json:"game_type"`
	Name       string `json:"name"`
	UserId     string `json:"user_id"`
}

type X19GameCharacterQueryResult struct {
	Code     int                           `json:"code"`
	Details  string                        `json:"details"`
	Entities []X19GameCharacterQueryEntity `json:"entities"`
	Message  string                        `json:"message"`
	Total    int                           `json:"total"`
}

type X19SingleCharacterResult struct {
	Code    int                         `json:"code"`
	Details string                      `json:"details"`
	Entity  X19GameCharacterQueryEntity `json:"entity"`
	Message string                      `json:"message"`
}

// user id will auto fill into query info
func X19GameCharacters(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, query X19GameCharacterQueryInfo, result *X19GameCharacterQueryResult) error {
	if query.UserId == "" {
		query.UserId = user.Id
	}

	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/game-character/query/user-game-characters", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

// user id will auto fill into query info
func X19CreateGameCharacter(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, query X19CreateGameCharacterInfo, result *X19SingleCharacterResult) error {
	if query.UserId == "" {
		query.UserId = user.Id
	}

	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/game-character", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

func X19DeleteGameCharacter(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, entityId string, deleteMode string, result *X19SingleCharacterResult) error {
	query := struct {
		EntityId string `json:"entity_id"`
	}{
		EntityId: entityId,
	}

	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/game-character/"+deleteMode, postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}
