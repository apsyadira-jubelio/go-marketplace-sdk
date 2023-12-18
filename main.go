package main

import (
	"log"
	"net/url"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tokopedia"
)

func main() {
	APIURL, _ := url.Parse("https://fs.tokopedia.net")

	appConfig := tokopedia.ProxyAppConfig{
		ProxyURL:  "",
		EnableLog: true,
		AppID:     16619,
		APIURL:    APIURL,
	}

	client := tokopedia.NewHTTPProxy(appConfig)

	options := client.SetAccessToken("").CreateParams("", "GET", nil)
	res, err := client.SendRequest("/api/proxy", options)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
