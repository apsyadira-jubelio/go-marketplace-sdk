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

	sendMsg, err := client.Chat.SendMessage(context.Background(), &lazada.SendMessageParams{
		SessionID:     "",
		TemplateID:    4,
		Txt:           "[cheer up]",
		ImgURLSticker: "https://sg-live.slatic.net/other/im/901116c5359491b3bfe57fdb299f893f.gif",
		SmallImgURL:   "https://sg-live.slatic.net/other/im/5e917d78922de317df3a25c2a3c13aa6.png",
	})
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(sendMsg)
}
