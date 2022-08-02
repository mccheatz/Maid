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
	c.GenerateMobile()
	app := rest.MPayAppInfo{}
	// app.GenerateForX19(session.LatestPatch)
	app.GenerateForX19Mobile("840204111")
	var device rest.MPayDevice
	err = rest.MPayDevices(&client, c, app, &device)
	if err != nil {
		panic(err)
	}

	var user rest.MPayUser
	// err = rest.MPayLogin(&client, device, app, c, "f1182916778@163.com", "020601", &user)
	// if err != nil {
	// 	panic(err)
	// }

	err = rest.MPayLoginGuest(&client, device, app, c, &user)
	if err != nil {
		panic(err)
	}

	if user.RealNameStatus == 0 { // not real-name verified
		println("attempt real-name verify...")
		var result rest.MPayRealNameResult
		err = rest.MPayRealNameUpdate(&client, device, app, c, user, "姓名", "86", "362321195502064333", &result)
		if err != nil {
			panic(err)
		}

		if result.RealNameType == "成年人" {
			println("real-name verified!")
		}
	}

	var otpEntity rest.X19OTPEntity
	err = rest.LoginOTP(&client, user.ConvertToSAuth("x19", c, device), session.UserAgent, session.Release.CoreServerUrl, &otpEntity)
	if err != nil {
		panic(err)
	}

	println("MPay UserToken: " + user.Token)
	println("X19 OTP_Token:" + otpEntity.Entity.OTPToken)
}
