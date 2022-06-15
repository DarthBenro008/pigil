package main

import (
	"gnoty/internal/database"
	service2 "gnoty/internal/service"
	"gnoty/internal/utils"
	"log"
)

func InsertCommand(service database.DatabaseService, information database.CommandInformation) {
	err := service.Insert(information)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ListCommand(service database.DatabaseService) {
	data, err := service.List()
	if err != nil {
		log.Fatal(err.Error())
	}
	utils.PrintInformation(data)
}

func GoogleAuth(service database.DatabaseService) {
	config := service2.OAuthGoogleConfig()
	service2.GoogleLogin(config)
}
