package main

import (
	"log"
)

func InsertCommand(service DatabaseService, information CommandInformation) {
	err := service.Insert(information)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ListCommand(service DatabaseService) {
	data, err := service.List()
	if err != nil {
		log.Fatal(err.Error())
	}
	PrintInformation(data)
}

func GoogleAuth(service DatabaseService) {
	config := OAuthGoogleConfig()
	GoogleLogin(config)
}
