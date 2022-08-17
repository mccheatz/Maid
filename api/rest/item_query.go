package rest

import (
	"encoding/json"
	"errors"
	"maid/util"
	"net/http"
	"strconv"
)

type X19ItemQueryInfo struct {
	ItemType     int `json:"item_type"`
	Length       int `json:"length"`
	MasterTypeId int `json:"master_type_id"`
	Offset       int `json:"offset"`
}

type X19ItemQueryEntity struct {
	AvailableScope  int    `json:"available_scope"`
	BalanceGrade    int    `json:"balance_grade"`
	BriefSummary    string `json:"brief_summary"`
	DeveloperName   string `json:"developer_name"`
	EffectMTypeId   int    `json:"effect_mtypeid"`
	EffectSTypeId   int    `json:"effect_stypeid"`
	EntityId        string `json:"entity_id"`
	GameStatus      int    `json:"game_status"`
	GoodsState      int    `json:"goods_state"`
	IsApollo        int    `json:"is_apollo"`
	IsAuth          bool   `json:"is_auth"`
	IsCurrentSeason bool   `json:"is_current_season"`
	IsHas           bool   `json:"is_has"`
	ItemType        int    `json:"item_type"`
	ItemVersion     string `json:"item_version"`
	LobbyMaxNum     int    `json:"lobby_max_num"`
	LobbyMinNum     int    `json:"lobby_min_num"`
	MasterTypeId    string `json:"master_type_id"`
	ModId           int    `json:"mod_id"`
	Name            string `json:"name"`
	OnlineCount     string `json:"online_count"`
	PublishTime     int64  `json:"publish_time"`
	RelId           string `json:"rel_iid"`
	ResourceVersion int    `json:"resource_version"`
	ReviewStatus    int    `json:"review_status"`
	SeasonBegin     int    `json:"season_begin"`
	SeasonNumber    int    `json:"season_number"`
	SecondaryTypeId string `json:"secondary_type_id"`
	VipOnly         bool   `json:"vip_only"`
}

type X19ItemQueryResult struct {
	Code     int                  `json:"code"`
	Details  string               `json:"details"`
	Entities []X19ItemQueryEntity `json:"entities"`
	Message  string               `json:"message"`
	Total    string               `json:"total"`
}
type X19ItemVersionQueryInfo struct {
	ItemId string `json:"item_id"`
	Length int    `json:"length"`
	Offset int    `json:"offset"`
}

type X19ItemVersionQueryEntity struct {
	EntityId    string `json:"entity_id"`
	ItemId      string `json:"item_id"`
	JavaVersion int    `json:"java_version"`
	MCVersionId string `json:"mc_version_id"`
}

func (entity X19ItemVersionQueryEntity) GetMcVersionCode() int {
	var result int
	switch entity.MCVersionId {
	case "1":
		result = 1007010
	case "2":
		result = 1008000
	case "3":
		result = 1009004
	case "5":
		result = 1011002
	case "6":
		result = 1008008
	case "7":
		result = 1010002
	case "8":
		result = 1006004
	case "9":
		result = 1007002
	case "10":
		result = 1012002
	// case "11": result = 0
	case "12":
		result = 1008009
	case "13":
		result = 100000000
	case "14":
		result = 1013002
	case "15":
		result = 1014003
	case "16":
		result = 1015000
	case "17":
		result = 1016000
	case "18":
		result = 20000000
	case "19":
		result = 1018000
	default:
		result = 0
	}

	return result
}

type X19ItemVersionQueryResult struct {
	Code     int                         `json:"code"`
	Details  string                      `json:"details"`
	Entities []X19ItemVersionQueryEntity `json:"entities"`
	Message  string                      `json:"message"`
	Total    int                         `json:"total"`
}

type X19ItemAddressQueryEntity struct {
	Announcement string `json:"announcement"`
	EntityId     string `json:"entity_id"`
	GameStatus   int    `json:"game_status"`
	InWhitelist  bool   `json:"ip_whitelist"`
	Ip           string `json:"ip"`
	ISPEnable    bool   `json:"isp_enable"`
	Port         int    `json:"port"`
}

type X19ItemAddressQueryResult struct {
	Code    int                       `json:"code"`
	Details string                    `json:"details"`
	Entity  X19ItemAddressQueryEntity `json:"entity"`
	Message string                    `json:"message"`
}

func X19ItemQuery(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, query X19ItemQueryInfo, result *X19ItemQueryResult) error {
	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/item/query/available", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

func X19FetchAllQuery(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, query X19ItemQueryInfo, entities *[]X19ItemQueryEntity) error {
	current := 0
	var result X19ItemQueryResult

	*entities = make([]X19ItemQueryEntity, 0)

	for {
		query.Offset = current
		err := X19ItemQuery(client, userAgent, user, release, query, &result)
		if err != nil {
			return err
		}

		*entities = append(*entities, result.Entities...)

		current += query.Length
		max, err := strconv.Atoi(result.Total)
		if err != nil {
			return err
		}
		if current >= max {
			break
		}
	}

	return nil
}

func X19ItemVersionQuery(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, query X19ItemVersionQueryInfo, result *X19ItemVersionQueryResult) error {
	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/item-mc-version/query/search-by-item", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

func X19ItemVersionQueryById(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, entityId string, result *X19ItemVersionQueryEntity) error {
	query := X19ItemVersionQueryInfo{
		ItemId: entityId,
		Length: 1,
		Offset: 0,
	}

	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/item-mc-version/query/search-by-item", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	var response X19ItemVersionQueryResult

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if len(response.Entities) == 0 {
		return errors.New("no entity found in response")
	}

	*result = response.Entities[0]

	return nil
}

func X19ItemAddress(client *http.Client, userAgent string, user util.X19User, release X19ReleaseInfo, itemId string, result *X19ItemAddressQueryResult) error {
	query := struct {
		ItemId string `json:"item_id"`
	}{
		ItemId: itemId,
	}

	postBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	body, err := util.X19SimpleRequest("POST", release.ApiGatewayUrl+"/item-address/get", postBody, client, userAgent, user)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}
