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
		PartnerID:  2005794,
		PartnerKey: "6971596a5361646e446358774557784e4452436b575057706754534359637648",
		APIURL:     apiURL,
		MaxTimeout: 15 * time.Second,
	}

	relPath := fmt.Sprintf("/sellerchat/get_one_conversation?conversation_id=%d", 195209425697517370)
	client := shopee.NewProxyClient(app)
	params := client.WithShopID(45449350, "704e4d586a6350774c526d6565455078").CreateParams(relPath, "GET", nil, "")
	spew.Dump(params)
	resp, err := client.SendRequest(params)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
