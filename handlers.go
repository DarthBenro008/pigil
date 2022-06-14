package main

import (
	"fmt"
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
	fmt.Println(data)
}
