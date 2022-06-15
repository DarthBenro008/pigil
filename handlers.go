package main

import (
	"fmt"
	"gnoty/internal/database"
	service2 "gnoty/internal/service"
	"gnoty/internal/types"
	"gnoty/internal/utils"
	"log"
)

func InsertCommand(service database.Service,
	information types.CommandInformation) {
	err := service.Insert(information)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ListCommand(service database.Service) {
	data, err := service.List()
	if err != nil {
		log.Fatal(err.Error())
	}
	utils.PrintInformation(data)
}

func GoogleAuth(service database.Service) {
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

	err := service.InsertConfig(userEmail)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = service.InsertConfig(userRT)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = service.InsertConfig(userAT)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Status(service database.Service) {
	list, err := service.ListConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(list)
}

func Logout(service database.Service) {
	err := service.DeleteConfigDb()
	if err != nil {
		log.Fatal(err.Error())
	}
}
