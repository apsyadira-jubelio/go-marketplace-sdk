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
		ProxyURL:   "https://sp-proxy.jubelio.com",
		PartnerID:  123,
		PartnerKey: "",
		APIURL:     apiURL,
		MaxTimeout: 15 * time.Second,
	}

	relPath := fmt.Sprintf("/sellerchat/get_one_conversation?conversation_id=%d", 13)
	client := shopee.NewProxyClient(app)
	params := client.WithShopID(1, "704e4d586a6350774c526d6565455078").CreateParams(relPath, "GET", nil, "")
	spew.Dump(params)
	resp, err := client.SendRequest(params)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
