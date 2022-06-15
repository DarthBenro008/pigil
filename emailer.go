package main

import (
	"context"
	"encoding/base64"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func SendEmail(client *http.Client, email string) {
	srv, err := gmail.NewService(context.Background(),
		option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err.Error())
	}
	emailBody := "This is an test email from gnoty!"
	if err != nil {
		log.Fatal(err.Error())
	}

	var message gmail.Message

	emailTo := "To: " + email + "\r\n"
	subject := "Subject: " + "Gnoty sends regards!" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Fatal(err.Error())
	}
}
