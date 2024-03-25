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
	resp, err := client.Auth.GetAuthorizationShop("ROW_Tm5R_QAAAABftY_-lBYbKUNezeTwBEzVkxf1Ds1UuyUNmR0NVerLbyWyVo1aqiEvJzoo7CDU6icxj_y-36qDAg3oZ01l2KtMXD0cOqnJdu93Q_WiBUwjt1NYYaw0ptvTivsbZ2gNw5Xk2qyyEcnCjX2nVnO2wfrbVgTUFco2Y9XYdcjEwuMkzw", "")
	if err != nil {
		log.Fatal(err.Error())
	}

	spew.Dump(resp)

}
