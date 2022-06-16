package main

import (
	"context"
	"fmt"
	"github.com/DarthBenro008/pigil/internal/database"
	service2 "github.com/DarthBenro008/pigil/internal/service"
	"github.com/DarthBenro008/pigil/internal/types"
	"github.com/DarthBenro008/pigil/internal/utils"
	"golang.org/x/oauth2"
)

const handlerTag = "handlers"

func InsertCommand(service database.Service,
	information types.CommandInformation) {
	err := service.Insert(information)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
}

func ListCommand(service database.Service) {
	data, err := service.List()
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	if len(*data) == 0 {
		utils.InformationLogger("No history found yet!")
		return
	}
	utils.PrintInformation(data)
}

func GoogleAuth(service database.Service) {
	email, err := service.GetConfig(utils.UserEmail)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	if email != "" {
		utils.InformationLogger(fmt.Sprintf("already logged in with %s\n",
			email))
		return
	}
	config := service2.OAuthGoogleConfig()
	config = service2.GoogleLogin(config)
	data := service2.GoogleCallback(config)
	userEmail := types.ConfigurationInformation{
		Key:   utils.UserEmail,
		Value: data.Email,
	}
	userRT := types.ConfigurationInformation{
		Key:   utils.UserAT,
		Value: data.AccessToken,
	}
	userAT := types.ConfigurationInformation{
		Key:   utils.UserRT,
		Value: data.RefreshToken,
	}

	err = service.InsertConfig(userEmail)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	err = service.InsertConfig(userRT)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	err = service.InsertConfig(userAT)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
}

func Status(service database.Service) {
	list, err := service.ListConfig()
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	if len(*list) == 0 {
		utils.InformationLogger("No settings found yet!")
		return
	}
	utils.StatusPrinter(list)
}

func Logout(service database.Service) {
	err := service.DeleteConfigDb()
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
}

func Notify(service database.Service, information types.CommandInformation) {
	email, err := service.GetConfig(utils.UserEmail)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	if email == "" {
		utils.ErrorInformation("You are not authenticated! pigil cannot" +
			" notify via email, please run `pigil bumf auth`")
		return
	}
	at, err := service.GetConfig(utils.UserAT)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	rt, err := service.GetConfig(utils.UserRT)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	cred := oauth2.Token{AccessToken: at, RefreshToken: rt}
	client := service2.OAuthGoogleConfig().Client(context.Background(), &cred)
	service2.SendEmail(client, email, information)
}

func Help() {
	msg := utils.GenerateHelp(Version)
	fmt.Println(msg)
}
