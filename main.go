package main

import (
	"fmt"
	"maid/api"
	"maid/api/rest"
	"maid/util"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	rand.Seed(time.Now().UnixMilli()) // reset random seed

	var err error

	var dial func(network, addr string) (c net.Conn, err error)
	{
		baseDialer := &net.Dialer{
			Timeout: 5 * time.Second,
		}
		// detect system proxy settings
		sysproxy := os.Getenv("SOCKS_PROXY")
		if sysproxy == "" {
			sysproxy = os.Getenv("SOCKS5_PROXY")
		}
		if sysproxy != "" {
			println("Proxy detected: " + sysproxy)
			dialProxy, _ := proxy.SOCKS5("tcp", sysproxy, nil, baseDialer)
			dial = dialProxy.Dial
		} else {
			dial = baseDialer.Dial
		}
	}
	client := &http.Client{
		Transport: &http.Transport{Dial: dial},
		Timeout:   5 * time.Second,
	}

	session, err := api.EstablishSession(client)
	if err != nil {
		panic(err)
	}
	err = session.CheckSessionAbility()
	if err != nil {
		panic(err)
	}

	clientMPay := rest.MPayClientInfo{}
	clientMPay.GeneratePC()
	clientMPay.Udid = "o0Oooo0oO"
	app := rest.MPayAppInfo{}
	// app.GenerateForX19(session.LatestPatch.Version)
	app.GenerateForX19Mobile(session.LatestPatch.Version)
	var device rest.MPayDevice
	err = rest.MPayDevices(client, clientMPay, app, &device)
	if err != nil {
		panic(err)
	}

	var user rest.MPayUser
	// err = rest.MPayLogin(client, device, app, clientMPay, "f1182916778@163.com", "020601", &user)
	err = rest.MPayLoginGuest(client, device, app, clientMPay, &user)
	if err != nil {
		panic(err)
	}

	println("MPay UserToken: " + user.Token)

	if user.RealNameStatus == 0 { // not real-name verified
		println("attempt real-name verify...")
		var result rest.MPayRealNameResult
		err = rest.MPayRealNameUpdate(client, device, app, user, "姓名", "86", "362321195502064333", &result)
		if err != nil {
			panic(err)
		}

		if result.RealNameType == "成年人" {
			println("real-name verified!")
		}
	}

	sAuth := user.ConvertToSAuth("x19", clientMPay, device)

	var otpEntity rest.X19OTPEntity
	err = rest.X19LoginOTP(client, sAuth, session.UserAgent, session.Release, &otpEntity)
	if err != nil {
		panic(err)
	}

	var authEntity rest.X19AuthenticationEntity
	err = rest.X19AuthenticationOTP(client, session.UserAgent, sAuth, clientMPay, session.Release, session.LatestPatch, otpEntity, &authEntity)
	if err != nil {
		panic(err)
	}

	fmt.Printf("X19 Auth (Token=%s, EntityId=%s)\n", authEntity.Entity.Token, authEntity.Entity.EntityId)

	x19User := authEntity.ToUser()

	// println("attempt establish connection to authenticate server...")
	// authConn := api.X19AuthServerConnection{
	// 	Address:   session.AuthServers[rand.Intn(len(session.AuthServers))].ToAddr(),
	// 	Dial:      dial,
	// 	UserToken: x19User.Token,
	// 	EntityId:  x19User.Id,
	// }
	// go func() {
	// 	err := authConn.Establish()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()

	// str := api.GenerateAuthenticationBody("3x18ae9a1bddabf26174d8941ed6159b01905131", "7711451783364710", "1.12.2", "1.8.24.54179", "MODS", "56ad4607b81cd5a5f0c829db665d5da1", "c4a310031cde5126bd524f8d403df764")

	// for !authConn.HasEstablished() {
	// 	time.Sleep(1 * time.Second)
	// }

	// println("connected!")
	// err = authConn.SendPacket(9, str)
	// if err != nil {
	// 	panic(err)
	// }

	// return

	// update session every minute is required
	updater := api.X19SessionUpdater{}
	go updater.StartUpdating(client, session.UserAgent, session.Release, &x19User)

	// fetch server list
	var serverItem rest.X19ItemQueryEntity
	{
		itemQuery := rest.X19ItemQueryInfo{
			ItemType:     1,
			Length:       50,
			MasterTypeId: 2,
		}
		var queryResultItem []rest.X19ItemQueryEntity
		err = rest.X19FetchAllQuery(client, session.UserAgent, x19User, session.Release, itemQuery, &queryResultItem)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d servers online\n", len(queryResultItem))

		// find the server that I intend to join
		for _, item := range queryResultItem {
			if item.Name == "花雨庭" {
				serverItem = item
			}
		}

		if serverItem.Name == "" {
			println("Failed to find the target server")
			return
		}
	}

	var serverAddress rest.X19ItemAddressQueryEntity
	{
		var query rest.X19ItemAddressQueryResult
		err = rest.X19ItemAddress(client, session.UserAgent, x19User, session.Release, serverItem.EntityId, &query)
		if err != nil {
			panic(err)
		}
		serverAddress = query.Entity
	}

	fmt.Printf("server found! (name=%s, id=%s, address=%s)\n", serverItem.Name, serverItem.EntityId, serverAddress.Address())

	// fetch character list
	var characters []rest.X19GameCharacterQueryEntity
	{
		characterQuery := rest.X19GameCharacterQueryInfo{
			GameId:   serverItem.EntityId,
			GameType: 2,
			Length:   50,
			Offset:   0,
		}
		var queryResultCharacter rest.X19GameCharacterQueryResult
		err = rest.X19GameCharacters(client, session.UserAgent, x19User, session.Release, characterQuery, &queryResultCharacter)
		if err != nil {
			panic(err)
		}
		characters = append(characters, queryResultCharacter.Entities...)
		fmt.Printf("%d character(s) found\n", len(characters))
	}

	if len(characters) == 0 {
		println("no character found! attempt create")

		characterCreateQuery := rest.X19CreateGameCharacterInfo{
			GameId:   serverItem.EntityId,
			GameType: 2,
			Name:     "Taka_" + util.RandStringRunes(5),
		}
		var queryResultCharacter rest.X19SingleCharacterResult
		err = rest.X19CreateGameCharacter(client, session.UserAgent, x19User, session.Release, characterCreateQuery, &queryResultCharacter)
		if err != nil {
			panic(err)
		}
		characters = append(characters, queryResultCharacter.Entity)
	}

	for _, c := range characters {
		t := time.Unix(c.CreateTime, 0)
		println(c.EntityId + "\t" + c.Name + "\t" + t.Local().Format("2006-01-02 15:04:05"))
	}

	var versionInfo rest.X19ItemVersionQueryEntity
	err = rest.X19ItemVersionQueryById(client, session.UserAgent, x19User, session.Release, serverItem.EntityId, &versionInfo)
	if err != nil {
		panic(err)
	}

	return

	// download game
	{
		// downloads := make([]util.DownloadInfo, 0)

		var itemResult rest.X19UserItemResult
		err = rest.X19UserItemDownload(client, session.UserAgent, x19User, session.Release, serverItem.EntityId, &itemResult)
		if err != nil {
			panic(err)
		}

		if len(itemResult.Entity.SubEntities) == 0 {
			println("no game version found")
			return
		}

		mods, err := rest.FetchGameResourcesVerifyList(client, itemResult.Entity.SubEntities)
		if err != nil {
			panic(err)
		}

		println(mods)

		// for _, sub := range itemResult.Entity.SubEntities {
		// 	downloads = append(downloads, util.DownloadInfo{
		// 		Path: "./dl/" + sub.ResourceName,
		// 		Url:  sub.ResourceUrl,
		// 	})
		// }

		/*
			query := rest.X19AuthItemQuery{
				GameType:    2,
				McVersionId: versionInfo.GetMcVersionCode(),
			}
			var authItemResult rest.X19AuthItemResult
			err = rest.X19AuthItemSearch(client, session.UserAgent, x19User, session.Release, query, &authItemResult)
			if err != nil {
				panic(err)
			}

			itemIds := make([]string, len(authItemResult.Entity.IIdList))
			for i, item := range authItemResult.Entity.IIdList {
				itemIds[i] = item.Value
			}

			var itemListResult rest.X19UserItemListResult
			err = rest.X19UserItemListDownload(client, session.UserAgent, x19User, session.Release, itemIds, &itemListResult)
			if err != nil {
				panic(err)
			}
		*/

		// for _, item := range itemListResult.Entities {
		// 	for _, sub := range item.SubEntities {
		// 		downloads = append(downloads, util.DownloadInfo{
		// 			Path: "./dl/" + sub.ResourceName,
		// 			Url:  sub.ResourceUrl,
		// 		})
		// 	}
		// }

		searchKeysQuery := rest.X19SearchKeysQuery{
			ForgeVersion:    versionInfo.GetMcVersionCode(),
			GameType:        2,
			ItemIdList:      make([]string, 0),
			ItemVersionList: make([]string, 0),
			ItemMd5List:     make([]string, 0),
		}
		var searchKeysResult rest.X19SearchKeysResult
		err = rest.X19SearchKeysByItemList(client, session.UserAgent, x19User, session.Release, searchKeysQuery, &searchKeysResult)
		if err != nil {
			panic(err)
		}

		var launchWrapperMD5 string
		var gameDataMD5 string
		for _, key := range searchKeysResult.Entities {
			if strings.HasSuffix(key.Name, ".dat") {
				gameDataMD5 = key.MD5
			} else if strings.Contains(key.Name, "launchwrapper") {
				launchWrapperMD5 = key.MD5
			}
		}

		if launchWrapperMD5 == "" || gameDataMD5 == "" {
			println("failed to fetch file integrity")
			return
		}

		println(launchWrapperMD5)
		println(gameDataMD5)
		println(session.LatestPatch.Version)

		// now := time.Now()
		// println("Downloading...")
		// errs := util.ParallelDownload(downloads, client)
		// println("downloaded! " + time.Since(now).String())
		// for _, e := range errs {
		// 	println(e.Error())
		// }
	}
}
