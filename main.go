package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/shopee"
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

	conversation, err := client.Voucher.GetListVoucherByStatus(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetVoucherListParam{
		PageNo:   1,
		PageSize: 25,
		Status:   "ongoing",
	})
	if err != nil {
		log.Fatal(err)
		return // Consider handling errors differently if partial results are acceptable
	}

	spew.Dump(conversation)
	writeJSONFile(conversation, "voucher-list")
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

	fmt.Println("Response written to response.json")
	return nil
}
