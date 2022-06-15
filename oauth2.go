package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

func OAuthGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:6969",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo." +
			"profile", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/gmail.compose"},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(config *oauth2.Config) {
	//TODO: state generation
	url := config.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Printf("Click on this link to authenticate yourself with gnoty! \n%s"+
		"\n", url)
	GoogleCallback(config)
}

func GoogleCallback(config *oauth2.Config) {
	server := http.Server{Addr: ":6969", Handler: nil}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			log.Fatal("could not rec", err.Error())
		}
		fmt.Println(token.AccessToken)
		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintf(w,
			"Your email has been linked via gnoty! You can close this webpage"+
				" now!")
		if err != nil {
			log.Fatal(err.Error())
		}
		SendEmail(config.Client(context.Background(), token),
			"") //TODO: Fetch user data
		err = server.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
