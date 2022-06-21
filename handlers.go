package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/DarthBenro008/pigil/internal/database"
	service2 "github.com/DarthBenro008/pigil/internal/service"
	"github.com/DarthBenro008/pigil/internal/types"
	"github.com/DarthBenro008/pigil/internal/utils"
	"golang.org/x/oauth2"
	"time"
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
	for _, channel := range utils.Channels {
		value, err := service.GetConfig(channel)
		if err != nil {
			utils.ErrorLogger(err, handlerTag)
		}
		if value == utils.DbTrue {
			switch channel {
			case utils.Email:
				err = emailNotifier(service, information)
				if err != nil {
					utils.ErrorLogger(err, handlerTag)
				}
			case utils.Discord:
				err = discordNotifier(service, information)
				if err != nil {
					utils.ErrorLogger(err, handlerTag)
				}
			}
		}
	}
}

func emailNotifier(service database.Service, information types.CommandInformation) error {
	email, err := service.GetConfig(utils.UserEmail)
	if err != nil {
		return err
	}
	if email == "" {
		utils.ErrorInformation("You are not authenticated! pigil cannot" +
			" notify via email, please run `pigil bumf auth`")
		return nil
	}
	at, err := service.GetConfig(utils.UserAT)
	if err != nil {
		return err
	}
	rt, err := service.GetConfig(utils.UserRT)
	if err != nil {
		return err
	}
	cred := oauth2.Token{AccessToken: at, RefreshToken: rt,
		Expiry: time.Now().Add(5)}
	client := service2.OAuthGoogleConfig().Client(context.Background(), &cred)
	service2.SendEmail(client, email, information)
	return nil
}

func discordNotifier(service database.Service, information types.CommandInformation) error {
	url, err := service.GetConfig(utils.DbDiscordUrl)
	msg := utils.GenerateDiscordMessage(information)
	err, status, resp := service2.SendPostRequest(url, msg)
	if status != 204 {
		return errors.New(resp)
	}
	if err != nil {
		return err
	}
	return nil
}

func NotificationSelector(service database.Service) {
	var selectedChannels []string
	prompt := &survey.MultiSelect{
		Message: "Select Channel of Notification",
		Options: utils.Channels,
		Default: utils.Email,
	}
	err := survey.AskOne(prompt, &selectedChannels, survey.WithValidator(func(ans interface{}) error {
		if len(ans.([]core.OptionAnswer)) == 0 {
			return errors.New("please select one of the following")
		}
		return nil
	}))
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	for _, channel := range utils.Channels {
		if utils.Contains(selectedChannels, channel) {
			err := service.InsertConfig(types.ConfigurationInformation{
				Key:   channel,
				Value: utils.DbTrue,
			})
			if err != nil {
				utils.ErrorLogger(err, handlerTag)
			}
		} else {
			err := service.InsertConfig(types.ConfigurationInformation{
				Key:   channel,
				Value: utils.DbFalse,
			})
			if err != nil {
				utils.ErrorLogger(err, handlerTag)
			}
		}
	}
}

func DiscordToggle(service database.Service, toggle bool) {
	if !toggle {
		err := service.InsertConfig(types.ConfigurationInformation{
			Key:   utils.Discord,
			Value: utils.DbTrue,
		})
		if err != nil {
			utils.ErrorLogger(err, handlerTag)
		}
		utils.InformationLogger("Discord Webhook has been disabled!")
		return
	}
	input := ""
	prompt := &survey.Input{
		Message: "Please enter discord webhook url: ",
	}
	err := survey.AskOne(prompt, &input, survey.WithValidator(survey.Required))
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	err = service.InsertConfig(types.ConfigurationInformation{
		Key:   utils.DbDiscordUrl,
		Value: input,
	})
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	err = service.InsertConfig(types.ConfigurationInformation{
		Key:   utils.Discord,
		Value: utils.DbTrue,
	})
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
}

func Help() {
	msg := utils.GenerateHelp(Version)
	fmt.Println(msg)
}

func IsFirstTime(service database.Service) {
	value, err := service.GetConfig(utils.DbFirstTime)
	if err != nil {
		utils.ErrorLogger(err, handlerTag)
	}
	if value == "" {
		fmt.Println("Welcome to Pigil! Please set the notification channels")
		NotificationSelector(service)
		err = service.InsertConfig(types.ConfigurationInformation{
			Key:   utils.DbFirstTime,
			Value: time.Now().String(),
		})
		if err != nil {
			utils.ErrorLogger(err, handlerTag)
		}
	}
}
