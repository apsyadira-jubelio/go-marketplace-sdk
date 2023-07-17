package main

import (
	"fmt"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/shopee"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := shopee.AppConfig{
		PartnerID:   2005794,
		PartnerKey:  "6971596a5361646e446358774557784e4452436b575057706754534359637648",
		RedirectURL: "",
		APIURL:      "https://partner.shopeemobile.com",
	}

	client := shopee.NewClient(app)
	res, err := client.Chat.UploadImage(45449350, "454f75704656506a516665694d614c61", "https://radarlampung.disway.id/upload/891504aea3381619b7bbf4670f20b785.jpg")

	if err != nil {
		spew.Dump(err)
	}

	fmt.Println("Success", res)
}
