package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/shopee"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// playground
	partnerID, _ := strconv.Atoi(os.Getenv("SHOPEE_PARTNER_ID"))
	shopID, _ := strconv.Atoi(os.Getenv("SHOP_ID"))
	APIURL, _ := url.Parse("https://partner.shopeemobile.com")

	appConfig := shopee.AppConfig{
		PartnerID:    partnerID,
		PartnerKey:   os.Getenv("SHOPEE_PARTNER_KEY"),
		RedirectURL:  "",
		APIURL:       APIURL.String(),
		EnableSocks5: true,
		SockAddress:  os.Getenv("SOCKS_ADDRESS"),
	}

	spew.Dump(appConfig)
	client := shopee.NewClient(appConfig, shopee.WithRetry(3), shopee.WithSocks5(os.Getenv("SOCKS_ADDRESS")))
	arrOrderID := make([]string, 0)
	batch, err := client.Order.GetOrderDetailByOrderSN(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetOrderDetailParamsRequest{
		OrderSNList: strings.Join(arrOrderID, ","),
	})

	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(batch)
}
