package main

import (
	"encoding/json"
	"fmt"
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
	client := shopee.NewClient(appConfig, shopee.WithRetry(3), shopee.WithSocks5(os.Getenv("SOCKS_ADDRESS")))
	// resp, err := client.Product.GetProductlList(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetProductListParamRequest{
	// 	PageSize:   100,
	// 	Offset:     0,
	// 	ItemStatus: "NORMAL",
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// itemsID := make([]int, 0, resp.Response.TotalCount)
	// for _, item := range resp.Response.Item {
	// 	itemsID = append(itemsID, int(item.ItemID))
	// }

	// Function to batch get product details
	// getProductDetailsInBatches := func(ids []int) ([]shopee.ItemListData, error) {
	// 	var products []shopee.ItemListData
	// 	batchSize := 50

	// 	for i := 0; i < len(ids); i += batchSize {
	// 		end := i + batchSize
	// 		if end > len(ids) {
	// 			end = len(ids)
	// 		}
	// 		batch, err := client.Product.GetProductById(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetProductParamRequest{
	// 			ItemIDList: ids[i:end],
	// 		})
	// 		if err != nil {
	// 			return nil, err // Consider handling errors differently if partial results are acceptable
	// 		}
	// 		products = append(products, batch.Response.ItemList...)
	// 	}

	// 	return products, nil
	// }

	conversation, err := client.Chat.GetMessage(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetMessageParamsRequest{
		Offset:         "0",
		PageSize:       50,
		ConversationID: 0,
	})
	if err != nil {
		log.Fatal(err)
		return // Consider handling errors differently if partial results are acceptable
	}

	spew.Dump(conversation)

	writeJSONFile(conversation, "get-messages")
}

func writeJSONFile(response interface{}, filename string) error {
	// Create a new JSON file
	file, err := os.Create(fmt.Sprintf("%s.json", filename))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Encode response data to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Indent the JSON for readability (optional)
	err = encoder.Encode(response)
	if err != nil {
		fmt.Println("Error encoding JSON to file:", err)
		return err
	}

	fmt.Printf("Response written to %s.json", filename)
	return nil
}
