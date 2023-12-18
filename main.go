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
		ProxyURL:   "",
		PartnerID:  123,
		PartnerKey: "",
		APIURL:     apiURL,
	}

	relPath := fmt.Sprintf("/product/get_item_base_info?item_id_list=%d&need_tax_info=true&need_complaint_policy=true", 123)
	client := shopee.NewProxyClient(app)
	params := client.WithShopID(123, "123").CreateParams(relPath, "GET", nil, "")
	resp, err := client.SendRequest(params)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
