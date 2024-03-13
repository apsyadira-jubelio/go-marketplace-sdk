package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/shopee"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	// playground
	APIURL, _ := url.Parse("https://partner.shopeemobile.com")
	appConfig := shopee.ProxyAppConfig{
		ProxyURL:    "https://sp-proxy.jubelio.com",
		PartnerID:   uint64(123),
		PartnerKey:  "",
		MaxTimeout:  1 * time.Minute,
		APIURL:      APIURL,
		AccessToken: "",
		ShopID:      123,
	}

	client := shopee.NewProxyClient(appConfig)
	// relPath := "/product/get_item_list?offset=0&page_size=5&item_status=NORMAL" // List without search
	relPath := "/product/search_item?item_name=indomie&page_size=1&item_status=NORMAL"
	params := client.CreateParams(relPath, "GET", nil, "")
	params.JSON = true

	spew.Dump(params)
	response, err := client.SendRequest(params)

	if err != nil {
		fmt.Println("error cause:", err)
		os.Exit(1)
	}
	var resp = make(map[string]interface{}, 0)
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		fmt.Println("error cause:", err)
		os.Exit(1)
	}

	if response.StatusCode() != 200 {
		fmt.Println("response not OK:", response.StatusCode())
		os.Exit(1)
	}

	spew.Dump(resp)
}
