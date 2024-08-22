package main

import (
	"log"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/lazada"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	appKey := ""
	appSecret := ""
	// playground
	client := lazada.NewClient(appKey, appSecret, lazada.Indonesia)
	resp, err := client.Chat.GetListSticker()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
}
