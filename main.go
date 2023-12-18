package main

import (
	"fmt"
	"log"
	"net/url"

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
	}

	relPath := fmt.Sprintf("/product/get_model_list?item_id=%d", 6308520269)
	client := shopee.NewProxyClient(app)
	params := client.WithShopID(45449350, "444344506646784b7573525279574178").CreateParams(relPath, "GET", nil, "")
	spew.Dump(params)
	resp, err := client.SendRequest(params)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
