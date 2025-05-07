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
	client := shopee.NewClient(appConfig, shopee.WithRetry(3), shopee.WithSocks5(os.Getenv("SOCKS_ADDRESS")))
	conversation, err := client.Chat.GetOneConversation(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.GetMessageParamsRequest{
		// Offset: "0",
		// PageSize: 5,
		ConversationID: 419812330742606495,
	})
	if err != nil {
		log.Fatal(err)
		return // Consider handling errors differently if partial results are acceptable
	}
	spew.Dump(conversation)

	readMessage, err := client.Chat.ReadConversation(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.ReadMessageRequest{
		ConversationID:    419812330742606495,
		BusinessType:      0,
		LastReadMessageID: "-",
	})
	if err != nil {
		log.Println("error while req to readMessage:", err)
	}

	spew.Dump(readMessage)
	unreadMessage, err := client.Chat.UnreadConversation(uint64(shopID), os.Getenv("SHOPEE_TOKEN"), shopee.UnreadMessageRequest{
		ConversationID: 419812330742606495,
		BusinessType:   0,
	})
	if err != nil {
		log.Println("error while req to unreadMessage:", err)
	}

	spew.Dump(unreadMessage)
}
