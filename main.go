package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/shopee"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	apiURL, _ := url.Parse("https://partner.shopeemobile.com")

	app := shopee.ProxyAppConfig{
		ProxyURL:   "",
		PartnerID:  123,
		PartnerKey: "123",
		APIURL:     apiURL,
		MaxTimeout: 15 * time.Second,
	}

	relPath := fmt.Sprintf("/product/get_model_list?item_id=%d", 123)
	client := shopee.NewProxyClient(app)
	params := client.WithShopID(123, "123").CreateParams(relPath, "GET", nil, "")
	spew.Dump(params)
	resp, err := client.SendRequest(params)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
