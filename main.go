package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tiktok"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// playground
	appConfig := tiktok.AppConfig{
		AppKey:    os.Getenv("TIKTOK_APP_KEY"),
		AppSecret: os.Getenv("TIKTOK_APP_SECRET"),
		APIURL:    "https://auth.tiktok-shops.com",
		Version:   "202309",
	}

	client := tiktok.NewClient(appConfig)

	state := map[string]string{
		"companyId": "123",
		"tenantId":  "1231212312",
		"mp":        "Tiktok",
		"platform":  "MACOS",
	}

	stateBytes, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err.Error()) // log.
	}

	stateBase64 := base64.StdEncoding.EncodeToString(stateBytes)
	url, err := client.Auth.GetLegacyAuthURL(appConfig.AppKey, stateBase64)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(url)
	qs := tiktok.GetAccessTokenParams{
		Code:      "123",
		AppKey:    appConfig.AppKey,
		AppSecret: appConfig.AppSecret,
		GrantType: "authorized_code",
	}

	resp, err := client.Auth.GetAccessToken(qs)
	if err != nil {
		log.Fatal(err.Error())
	}

	spew.Dump(resp)
}
