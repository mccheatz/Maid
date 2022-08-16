package main

import (
	"maid/api"
	"maid/api/rest"
	"maid/util"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixMilli()) // reset random seed

	// a := util.ComputeDynamicToken("/item-buy-list/query/search-mcgame-item-list-v2", []byte(`{"versions":["1.7.10","1.8","1.8.8","1.8.9","1.9.4","1.10.2","1.11.2","1.12.2","1.13.2","1.14.3","1.15","1.16","1.18"],"offset":0,"length":50}`), "")
	a := util.ComputeDynamicToken("/user-detail/89602", make([]byte, 0), "y9LkDerv903ECe6M")
	println(a)

	return

	proxyUrl, _ := url.Parse("http://127.0.0.1:8889")
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
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
	// app.GenerateForX19(session.LatestPatch)
	app.GenerateForX19Mobile("840204111")
	var device rest.MPayDevice
	err = rest.MPayDevices(client, clientMPay, app, &device)
	if err != nil {
		panic(err)
	}

	var user rest.MPayUser
	// err = rest.MPayLogin(client, device, app, c, "f1182916778@163.com", "020601", &user)
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

	// x19User := authEntity.ToUser()

	// update session every minute is required
	// err = rest.X19AuthenticationUpdate(client, session.UserAgent, session.Release, x19User)
	// if err != nil {
	// 	panic(err)
	// }

	println("X19 AuthToken: " + authEntity.Entity.Token)
}
