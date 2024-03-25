package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tiktok"
	"github.com/caarlos0/env"
	"github.com/jarcoal/httpmock"
	"github.com/joho/godotenv"
)

const (
	maxRetries  = 3
	shopID      = 1234567
	merchantID  = 0
	accessToken = "accesstoken"
)

var (
	client *tiktok.TiktokClient
	app    tiktok.AppConfig
)

func setup() {
	err := godotenv.Load()
	if err != nil {
		app = tiktok.AppConfig{
			AppKey:    os.Getenv("TIKTOK_APP_KEY"),
			AppSecret: os.Getenv("TIKTOK_APP_SECRET"),
			APIURL:    tiktok.OpenAPIURL,
			Version:   "202309",
		}
	} else {
		env.Parse(&app)
	}

	client = tiktok.NewClient(app,
		tiktok.WithRetry(maxRetries))
	httpmock.ActivateNonDefault(client.Client)
}

func teardown() {
	httpmock.DeactivateAndReset()
}

func loadFixture(filename string) []byte {
	f, err := ioutil.ReadFile("../../mockdata/tiktok/" + filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot load fixture %v", filename))
	}
	return f
}

func loadMockData(filename string, out interface{}) {
	f, err := ioutil.ReadFile("../../mockdata/tiktok/" + filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot load fixture %v", filename))
	}
	if err := json.Unmarshal(f, &out); err != nil {
		panic(fmt.Sprintf("decode mock data error: %s", err))
	}
}
