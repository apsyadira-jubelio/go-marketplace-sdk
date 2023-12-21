package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tokopedia"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	APIURL, _ := url.Parse("https://fs.tokopedia.net")
	appConfig := tokopedia.ProxyAppConfig{
		ProxyURL: "",
		AppID:    123,
		APIURL:   APIURL,
	}

	client := tokopedia.NewProxyClient(appConfig)
	relPath := fmt.Sprintf("/inventory/v1/fs/%d/product/info?product_id=%d", 123, 123)
	params := client.SetAccessToken("").CreateParams(relPath, "GET", nil)
	response, err := client.SendRequest("/api/proxy", params)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(response.Body())
}
