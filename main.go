package main

import (
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
		APIURL:    tiktok.AuthBaseURL,
	}

	client := tiktok.NewClient(appConfig)
	url, err := client.Auth.GetAuthURL(os.Getenv("TIKTOK_SERVICE_ID"))
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(appConfig)
	spew.Dump(url)

	resp, err := client.Auth.GetAccessToken("12345", "authorized_code")
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
}
