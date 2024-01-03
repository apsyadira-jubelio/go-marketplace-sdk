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
	app := shopee.ProxyAppConfig{
		ProxyURL:   "",
		PartnerID:  uint64(1234),
		PartnerKey: "",
		APIURL:     APIURL,
		MaxTimeout: 1 * time.Minute,
	}

	client := shopee.NewProxyClient(app)

	relPath := "/auth/token/get"
	body := map[string]interface{}{
		"code":       "",
		"shop_id":    123,
		"partner_id": 123,
	}
	params := client.CreateParams(relPath, "POST", body, "")

	params.JSON = true

	spew.Dump(params)
	response, err := client.SendRequest(params)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if response.StatusCode() != 200 {
		fmt.Println(response)
		os.Exit(1)
	}

	var token shopee.AccessTokenResponse

	err = json.Unmarshal(response.Body(), &token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spew.Dump(token)
}
