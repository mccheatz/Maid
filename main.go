package main

import (
	"maid/api"
	"maid/api/rest"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix()) // reset random seed

	client := http.Client{}

	// var info rest.BoxBasicInformation
	// rest.InitBox(&client, &info)
	// println(info.EncodeUrl)

	session, err := api.EstablishSession(&client)
	if err != nil {
		panic(err)
	}
	err = session.CheckSessionAbility()
	if err != nil {
		panic(err)
	}

	c := rest.MPayClientInfo{}
	c.Generate()
	app := rest.MPayAppInfo{}
	app.GenerateForX19(session.LatestPatch)
	var device rest.MPayDevice
	err = rest.MPayDevices(&client, c, app, &device)
	if err != nil {
		panic(err)
	}

	var user rest.MPayUser
	err = rest.MPayLogin(&client, device, app, c, "f1182916778@163.com", "02061101", &user)
	if err != nil {
		panic(err)
	}

	// var otpEntity rest.X19OTPEntity
	// err = rest.LoginOTP(&client, user.ConvertToSAuth("x19", c, device), session.UserAgent, session.Release.CoreServerUrl, &otpEntity)
	// if err != nil {
	// 	panic(err)
	// }

	// println(user.Token)
	// println(otpEntity.Entity.OTPToken)
}
