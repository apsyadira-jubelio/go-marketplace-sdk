package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tiktok"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// appKey := ""
	// appSecret := ""
	// playground
	client := tiktok.NewClient(tiktok.AppConfig{
		AppKey:    "",
		AppSecret: "",
		APIURL:    tiktok.OpenAPIURL,
		Version:   "202309",
	}, tiktok.WithRetry(3))
	client.AccessToken = ""
	client.ShopCipher = ""
	client.ShopID = "7495786382084246013"
	resp, err := client.Promotion.SearchCoupons(5, "", tiktok.SearchCouponsBody{
		Status:       []string{"ONGOING"},
		TitleKeyword: "Coup",
	})
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(resp)
	writeJSONFile(resp, "test-search-page-token-1")
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
