package rest

import (
	"encoding/json"
	"maid/util"
	"net/http"
)

type X19UserItemResult struct {
	Code    int               `json:"code"`
	Details string            `json:"details"`
	Entity  X19UserItemEntity `json:"entity"`
	Message string            `json:"message"`
}

type X19UserItemListResult struct {
	Code     int                 `json:"code"`
	Details  string              `json:"details"`
	Entities []X19UserItemEntity `json:"entities"`
	Message  string              `json:"message"`
}

type X19UserItemEntity struct {
	DownloadTime int    `json:"download_time"`
	EntityId     string `json:"entity_id"`
	ItemId       string `json:"item_id"`
	IType        int    `json:"itype"`
	MTypeId      int    `json:"mtypeid"`
	STypeId      int    `json:"stypeid"`
	SubEntities  []struct {
		EntityId        string `json:"entity_id"`
		JarMD5          string `json:"jar_md5"`
		JavaVersion     int    `json:"java_version"`
		McVersionName   string `json:"mc_version_name"`
		ResourceMD5     string `json:"res_md5"`
		ResourceName    string `json:"res_name"`
		ResourceSize    int    `json:"res_size"`
		ResourceUrl     string `json:"res_url"`
		ResourceVersion int    `json:"res_version"`
	} `json:"sub_entities"`
	SubModList []util.JsonRaw `json:"sub_mod_list"`
	UserId     string         `json:"user_id"`
}

type X19AuthItemQuery struct {
	GameType    int `json:"game_type"`
	McVersionId int `json:"mc_version_id"`
}

type X19AuthItemResult struct {
	Code    int               `json:"code"`
	Details string            `json:"details"`
	Entity  X19AuthItemEntity `json:"entity"`
	Message string            `json:"message"`
}

type X19AuthItemEntity struct {
	GameType    int            `json:"game_type"`
	IIdList     []util.JsonRaw `json:"iid_list"`
	McVersionId int            `json:"mc_version_id"`
}

type X19SearchKeysQuery struct {
	ForgeVersion    int      `json:"forge_version"`
	ItemIdList      []string `json:"item_id_list"`
	ItemVersionList []string `json:"item_version_list"`
	ItemMd5List     []string `json:"item_md5_list"`
	GameType        int      `json:"game_type"`
	IsHost          int      `json:"is_host"`
}

type X19SearchKeysResult struct {
	Code     int                   `json:"code"`
	Details  string                `json:"details"`
	Entities []X19SearchKeysEntity `json:"entities"`
	Message  string                `json:"message"`
	Total    int                   `json:"total"`
}

type X19SearchKeysEntity struct {
	EntityId    string `json:"entity_id"`
	ItemId      string `json:"item_id"`
	ItemVersion string `json:"item_version"`
	Priority    int    `json:"priority"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	MD5         string `json:"md5"`
}

func X19UserItemDownload(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, itemId string, result *X19UserItemResult) error {
	query := struct {
		ItemId string `json:"item_id"`
	}{
		ItemId: itemId,
	}

	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/user-item-download-v2", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

func X19UserItemListDownload(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, items []string, result *X19UserItemListResult) error {
	query := struct {
		ItemIdList []string `json:"item_id_list"`
	}{
		ItemIdList: items,
	}

	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/user-item-download-v2/get-list", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

func X19AuthItemSearch(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, query X19AuthItemQuery, result *X19AuthItemResult) error {
	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/game-auth-item-list/query/search-by-game", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

func X19SearchKeysByItemList(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, query X19SearchKeysQuery, result *X19SearchKeysResult) error {
	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19EncryptRequest("POST", release.ApiGatewayUrl+"/item-key/query/search-keys-by-item-list-v2", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}
