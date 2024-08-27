package account

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type Account struct {
	service *oauth2.Service
}

func NewAccount() (*Account, error) {
	service, err := oauth2.NewService(context.Background(), option.WithCredentialsFile("private.json"))
	if err != nil {
		return nil, err
	}

	return &Account{
		service: service,
	}, nil
}

func (a *Account) GetUserInfo(accessToken string) {
	userInfo, err := a.service.Userinfo.Get().Do(googleapi.QueryParameter("access_token", accessToken))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(userInfo)
}
