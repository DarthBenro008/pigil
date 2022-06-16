package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/DarthBenro008/pigil/internal/types"
	"github.com/DarthBenro008/pigil/internal/utils"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"net/http"
)

const mailerTag = "mailer"

func SendEmail(client *http.Client, email string, information types.CommandInformation) {
	srv, err := gmail.NewService(context.Background(),
		option.WithHTTPClient(client))
	if err != nil {
		utils.ErrorLogger(err, mailerTag)
	}
	emailBody, err := utils.GenerateEmailTemplate(information)
	if err != nil {
		utils.ErrorLogger(err, mailerTag)
	}
	var message gmail.Message
	var status string
	if information.WasSuccessful {
		status = "passed"
	} else {
		status = "failed"
	}
	emailTo := fmt.Sprintf("To: %s\r\n", email)
	subject := fmt.Sprintf("Subject: [Pigil Notification]: Process `%s` %s\n",
		information.CommandName, status)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8" +
		"\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		utils.ErrorLogger(err, mailerTag)
	}
}
