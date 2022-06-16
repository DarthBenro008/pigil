package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DarthBenro008/pigil/internal/types"
	"github.com/DarthBenro008/pigil/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

const oauthTag = "oauth2"

func OAuthGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:6969",
		ClientID:     utils.GoogleClientId,
		ClientSecret: utils.GoogleClientSecret,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo." +
			"profile", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/gmail.compose"},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(config *oauth2.Config) *oauth2.Config {
	//TODO: state generation
	url := config.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Printf("Click on this link to authenticate yourself with github.com/DarthBenro008/pigil! \n%s"+
		"\n", url)
	return config
}

func GoogleCallback(config *oauth2.Config) types.UserInformation {
	server := http.Server{Addr: ":6969", Handler: nil}
	var userData types.UserInformation

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			utils.ErrorLogger(err, oauthTag)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			utils.ErrorLogger(err, oauthTag)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintf(w,
			"Your email has been linked via github.com/DarthBenro008/pigil! You can close this webpage"+
				" now!")
		if err != nil {
			utils.ErrorLogger(err, oauthTag)
		}
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			utils.ErrorLogger(err, oauthTag)
		}
		googleResponse := types.GoogleResponse{}
		err = json.NewDecoder(resp.Body).Decode(&googleResponse)
		if err != nil {
			utils.ErrorLogger(err, oauthTag)
		}
		userData = types.UserInformation{
			Name:         googleResponse.GivenName + " " + googleResponse.FamilyName,
			Email:        googleResponse.Email,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		go func() {
			err = server.Shutdown(context.Background())
			if err != nil {
				utils.ErrorLogger(err, oauthTag)
			}
		}()
	})
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		utils.ErrorLogger(err, oauthTag)
	}
	return userData
}
