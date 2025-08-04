package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/lazada"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	appKey := os.Getenv("LAZADA_APP_KEY")
	appSecret := os.Getenv("LAZADA_APP_SECRET")
	token := os.Getenv("LAZADA_TOKEN")
	// playground
	client := lazada.NewClient(appKey, appSecret, lazada.Indonesia)
	client.NewTokenClient(token)

	configLazada := map[string]string{
		"appKey":    appKey,
		"appSecret": appSecret,
		"token":     token,
	}

	spew.Dump(configLazada)

	resp, err := client.Voucher.GetVouchers(context.Background(), &lazada.GetVouchersParam{
		CurPage:     "1",
		VoucherType: "COLLECTIBLE_VOUCHER",
	})
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
	writeJSONFile(resp, "response-list-voucher")
}

func writeJSONFile(response any, filename string) error {
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
