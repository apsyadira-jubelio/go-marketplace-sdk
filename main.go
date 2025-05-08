package main

import (
	"log"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tiktok"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// appKey := ""
	// appSecret := ""
	// playground
	client := tiktok.NewClient(tiktok.AppConfig{
		AppKey:    "",
		AppSecret: "",
		APIURL:    tiktok.OpenAPIURL,
		Version:   "202309",
	}, tiktok.WithRetry(3))
	client.AccessToken = ""
	client.ShopCipher = ""
	client.ShopID = ""
	resp, err := client.Chat.ReadMessageConversationID("")
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
}
