package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tiktok"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// playground
	appConfig := tiktok.AppConfig{
		AppKey:    os.Getenv("TIKTOK_APP_KEY"),
		AppSecret: os.Getenv("TIKTOK_APP_SECRET"),
		APIURL:    tiktok.OpenAPIURL,
		Version:   "202309",
	}

	client := tiktok.NewClient(appConfig)

	state := map[string]string{
		"companyId": "123",
		"tenantId":  "1231212312",
		"mp":        "Tiktok",
		"platform":  "MACOS",
	}

	stateBytes, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err.Error()) // log.
	}

	stateBase64 := base64.StdEncoding.EncodeToString(stateBytes)
	url, err := client.Auth.GetLegacyAuthURL(appConfig.AppKey, stateBase64)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(url)

	respAccessToken, err := client.Auth.GetAccessToken(tiktok.GetAccessTokenParams{
		AppKey:    appConfig.AppKey,
		AppSecret: appConfig.AppSecret,
		Code:      "123",
		GrantType: "authorized_code",
	})
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(respAccessToken)

	respShop, err := client.Auth.GetAuthorizationShop("", "")
	if err != nil {
		log.Fatal(err.Error())
	}
	spew.Dump(respShop)
	shopCipher := ""
	shopID := ""
	for _, dataInShop := range respShop.Data.Shops {
		if dataInShop.Name == "JBBeauty" {
			shopCipher = dataInShop.Cipher
			shopID = dataInShop.ID
			break
		}
	}
	respRefreshToken, err := client.Auth.GetRefreshToken(tiktok.GetRefreshTokenParams{
		AppKey:       appConfig.AppKey,
		AppSecret:    appConfig.AppSecret,
		RefreshToken: respAccessToken.Data.RefreshToken,
		GrantType:    "refresh_token",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	spew.Dump(respRefreshToken)

	client.ShopCipher = shopCipher
	accessToken := "ROW_vYcihwAAAABftY_"
	convesationParam := tiktok.GetConversationsParam{
		PageSize: 2,
	}
	respConversations, err := client.Chat.GetConversations(convesationParam, shopCipher, shopID, accessToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	spew.Dump(respConversations)

	respConvMsg, err := client.Chat.GetConversationMessages(tiktok.GetConversationMessagesParam{
		PageSize: 3,
	}, "7320258856670609669", shopCipher, shopID, accessToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	spew.Dump(respConvMsg)
}
