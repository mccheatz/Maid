package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"maid/util"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bodgit/sevenzip"
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

type X19UserItemSubEntity struct {
	EntityId        string `json:"entity_id"`
	JarMD5          string `json:"jar_md5"`
	JavaVersion     int    `json:"java_version"`
	McVersionName   string `json:"mc_version_name"`
	ResourceMD5     string `json:"res_md5"`
	ResourceName    string `json:"res_name"`
	ResourceSize    int    `json:"res_size"`
	ResourceUrl     string `json:"res_url"`
	ResourceVersion int    `json:"res_version"`
}

type X19UserItemEntity struct {
	DownloadTime int                    `json:"download_time"`
	EntityId     string                 `json:"entity_id"`
	ItemId       string                 `json:"item_id"`
	IType        int                    `json:"itype"`
	MTypeId      int                    `json:"mtypeid"`
	STypeId      int                    `json:"stypeid"`
	SubEntities  []X19UserItemSubEntity `json:"sub_entities"`
	SubModList   []util.JsonRaw         `json:"sub_mod_list"`
	UserId       string                 `json:"user_id"`
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

type gameResourceVerify struct {
	ModPath string `json:"modPath"`
	Name    string `json:"name"`
	Id      string `json:"id"`
	ItemID  string `json:"iid"`
	MD5     string `json:"md5"`
}

func FetchGameResourcesVerifyList(client *http.Client, entities []X19UserItemSubEntity) (string, error) {
	checksum := ""
	for _, item := range entities {
		checksum += item.McVersionName + "-" + item.ResourceMD5 + ";"
	}
	checksum = "./.cache/" + util.MD5Hex([]byte(checksum)) + ".json"

	if _, err := os.Stat(checksum); err == nil {
		// path/to/whatever exists
		dat, err := os.ReadFile(checksum)
		if err != nil {
			return "", err
		}
		return string(dat), nil
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		entry := make([]gameResourceVerify, 0)

		for _, item := range entities {
			err := fetchGameResourcesVerify(client, item, &entry)
			if err != nil {
				return "", err
			}
		}

		jBody := struct {
			Mods []gameResourceVerify `json:"mods"`
		}{
			Mods: entry,
		}
		body, err := json.Marshal(jBody)
		if err != nil {
			return "", err
		}

		err = os.MkdirAll(path.Dir(checksum), os.ModePerm)
		if err != nil {
			return "", err
		}
		err = os.WriteFile(checksum, body, os.ModePerm)
		if err != nil {
			return "", err
		}

		return string(body), nil
	} else {
		return "", err
	}
}

func fetchGameResourcesVerify(client *http.Client, entity X19UserItemSubEntity, entry *[]gameResourceVerify) error {
	resp, err := client.Get(entity.ResourceUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	archive, err := sevenzip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	for _, file := range archive.File {
		if !file.FileInfo().IsDir() && strings.HasSuffix(file.Name, ".jar") && strings.Contains(file.Name, "/mods/") {
			reader, err := file.Open()
			if err != nil {
				return err
			}
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			reader.Close()
			dataMd5 := strings.ToUpper(util.MD5Hex(data))
			path := file.Name[strings.LastIndex(file.Name, "/")+1:]
			*entry = append(*entry, gameResourceVerify{
				ModPath: path,
				Id:      path,
				ItemID:  path[:strings.Index(path, "@")],
				MD5:     dataMd5,
			})
		}
	}

	return nil
}
