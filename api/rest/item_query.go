package rest

import (
	"encoding/json"
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
	PublishTime     int    `json:"publish_time"`
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
