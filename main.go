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
	initFileBytes, err := os.ReadFile("./test.mp4")
	if err != nil {
		log.Fatal(err)
	}

	uploadVideoResp, err := client.Chat.UploadVideo(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), "test.mp4", initFileBytes)
	if err != nil {
		log.Fatal(err)
		// Consider handling errors differently if partial results are acceptable
	}
	fmt.Println("Upload Video Response")
	spew.Dump(uploadVideoResp)

	respGetVideo, err := client.Chat.GetVideoByVidID(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetVideoParamRequest{
		Vid: uploadVideoResp.Response.Vid,
	})
	if err != nil {
		log.Fatal(err)
		// Consider handling errors differently if partial results are acceptable
	}
	fmt.Println("Get Video Response")
	spew.Dump(respGetVideo)

	resp, err := client.Chat.SendMessage(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.SendMessageRequest{
		MessageType:  "video",
		ToID:         json.Number(""),
		BusinessType: 0,
		Content: shopee.ContentSendMessage{
			Vid:      uploadVideoResp.Response.Vid,
			VideoURL: respGetVideo.Response.Video,
		},
	})
	if err != nil {
		log.Fatal(err)
		// Consider handling errors differently if partial results are acceptable
	}

	fmt.Println("Send Message Response")
	spew.Dump(resp)
	writeJSONFile(resp, "send-message resp")
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
