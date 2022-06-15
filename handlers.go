package main

import (
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
	service2.GoogleLogin(config)
}
