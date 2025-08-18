package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/lazada"
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

	// resp, err := client.Order.GetOrders(context.Background(), &lazada.GetOrdersParam{
	// 	CreatedAfter:  "2025-01-01T00:00:00Z",
	// 	SortBy:        "created_at",
	// 	SortDirection: "desc",
	// 	Limit:         "100",
	// 	Offset:        "0",
	// 	Status:        "packed",
	// })

	resp, err := client.Logistic.GetOrderTrace(context.Background(), &lazada.GetOrderTraceParams{
		OrderID: "2649549780800326",
	})

	if err != nil {
		log.Fatal(err)
	}
	writeJSONFile(resp, "response-list-order")
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
