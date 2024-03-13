package main

import (
	"log"
	"net/url"
	"os"
	"strconv"

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
	client := shopee.NewClient(appConfig, shopee.WithRetry(3))
	resp, err := client.Product.GetProductlList(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetProductListParamRequest{
		PageSize:   100,
		Offset:     0,
		ItemStatus: "NORMAL",
	})

	if err != nil {
		log.Fatal(err)
	}

	itemsID := make([]int, 0, resp.Response.TotalCount)
	for _, item := range resp.Response.Item {
		itemsID = append(itemsID, int(item.ItemID))
	}

	// Function to batch get product details
	getProductDetailsInBatches := func(ids []int) ([]shopee.ItemListData, error) {
		var products []shopee.ItemListData
		batchSize := 50

		for i := 0; i < len(ids); i += batchSize {
			end := i + batchSize
			if end > len(ids) {
				end = len(ids)
			}
			batch, err := client.Product.GetProductById(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetProductParamRequest{
				ItemIDList: ids[i:end],
			})
			if err != nil {
				return nil, err // Consider handling errors differently if partial results are acceptable
			}
			products = append(products, batch.Response.ItemList...)
		}

		return products, nil
	}

	products, err := getProductDetailsInBatches(itemsID)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(products)
}
