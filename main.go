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
		APIURL:    tiktok.OpenAPIURL,
		Version:   "202309",
	}

	client := tiktok.NewClient(appConfig)

	// use with common param request to set shop_id, shop_cipher, access_token
	// to access the OpenAPI
	client.WithCommonParamRequest(tiktok.CommonParamRequest{
		ShopID:      "123",
		ShopCipher:  "123",
		AccessToken: "123",
	})

	resp, err := client.Chat.ReadMessageConversationID("7350198843860353287")

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
