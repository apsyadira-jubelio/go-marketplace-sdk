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
	resp, err := client.Media.GetVideo(context.Background(), &lazada.GetVideoParameter{
		VideoID: "3026611400",
	})
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
}
