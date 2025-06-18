package main

import (
	"context"
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
	client.NewTokenClient("")
	resp, err := client.Chat.SendMessage(context.Background(), &lazada.SendMessageParams{
		SessionID:  "",
		TemplateID: lazada.EmojiMessage,
		Txt:        "[bye]",
		ImgUrl:     "https://sg-live.slatic.net/other/im/2a693f4624ae1491d6ba5743fd8a3ee9.gif",
	})
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
}
