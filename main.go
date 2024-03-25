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
		APIURL:    tiktok.OpenAPIURL,
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
	resp, err := client.Auth.GetAuthorizationShop("ROW_tZqEowAAAABftY_-lBYbKUNezeTwBEzV7T-uEdQR3qD7lu7tdl0YuX1OsYoBtH2L1nlzgH-m4OYORtNg3YKqUPBdiuleV17Tnndh8v9jpeM4Zk-pinJ7V19-G2fmQSgDu49cpezv52oc_aTopWJ-yClT2KmEKMZ7Mc-oLfpM4SizMW3CnXKG2g", "")
	if err != nil {
		log.Fatal(err.Error())
	}

	spew.Dump(resp)

}
