package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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
	now := time.Now().Local() // now, 21 march
	tanggal_7 := now.AddDate(0, 0, -14)
	resp, err := client.Order.GetListOrder(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetListOrderParamsRequest{
		PageSize:       100,
		TimeRangeField: "create_time",
		TimeFrom:       int(tanggal_7.Unix()),
		TimeTo:         int(now.Unix()),
		// OrderStatus:    "UNPAID",
	})

	if err != nil {
		log.Fatal(err)
	}

	if resp == nil || len(resp.Response.OrderList) == 0 {
		log.Fatal("order nil")
	}

	orderIDs := make([]string, 0, len(resp.Response.OrderList))
	for _, order := range resp.Response.OrderList {
		orderIDs = append(orderIDs, order.OrderSn)
	}

	getOrderDetailsInBatches := func(ids []string) ([]shopee.OrderList, error) {
		var orders []shopee.OrderList
		batchSize := 50

		for i := 0; i < len(ids); i += batchSize {
			end := i + batchSize
			if end > len(ids) {
				end = len(ids)
			}
			log.Println("check index array:", ids[i:end])
			batch, err := client.Order.GetOrderDetailByOrderSN(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetOrderDetailParamsRequest{
				OrderSNList: strings.Join(ids[i:end], ","),
			})
			if err != nil {
				log.Printf("error cause:%+v\n", err)
				return nil, err
			}
			// log.Printf("check data number %d we print %+v\n:", i, batch)
			orders = append(orders, batch.OrderListResponse.OrderList...)
		}

		return orders, nil
	}

	orders, err := getOrderDetailsInBatches(orderIDs)
	if err != nil {
		log.Printf("error cause:%+v\n", err)
		return
	}

	spew.Dump(orders)
}
