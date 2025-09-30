package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tiktok"
)

func main() {
	godotenv.Load()

	client := tiktok.NewClient(tiktok.AppConfig{
		AppKey:    "",
		AppSecret: "",
		APIURL:    tiktok.OpenAPIURL,
		Version:   "202309",
	}, tiktok.WithRetry(3))

	pageToken := "" // kosong di awal
	pageSize := 10

	allTokpedParticipants := []tiktok.Participants{}

	for {
		client.AccessToken = ""
		client.ShopCipher = ""
		client.ShopID = ""

		resp, err := client.Chat.GetConversations(tiktok.GetConversationsParam{
			PageSize:  pageSize,
			PageToken: pageToken,
		})
		if err != nil {
			log.Fatal("error get conversations:", err)
		}

		// filter buyer Tokopedia
		for _, eachData := range resp.Data.Conversations {
			for _, p := range eachData.Participants {
				if p.BuyerPlatform == "TOKOPEDIA" {
					allTokpedParticipants = append(allTokpedParticipants, p)
				}
			}
		}

		if resp.Data.NextPageToken == "" {
			break
		}
		pageToken = resp.Data.NextPageToken
	}

	err := writeJSONFile(allTokpedParticipants, "test-2")
	if err != nil {
		log.Println("Error while create JSON file:", err)
	}

}

func writeJSONFile(response interface{}, filename string) error {
	file, err := os.Create(fmt.Sprintf("%s.json", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}
