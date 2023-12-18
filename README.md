# go-marketplace-sdk

Marketplace SDK with Golang and Resty Client. Currently this SDK onlye support for Tokopedia, lazada, and shopee.

## How to use

### Shopee

Initialize Client And request shop info

```
  app := shopee.AppConfig{
		PartnerID:   123,
		PartnerKey:  "123",
		RedirectURL: "",
		APIURL:      "https://partner.shopeemobile.com",
	}

	shopeeClient := shopee.NewClient(app)

  authUrl, err := shopeeClient.Auth.GetAuthURL()
	if err != nil {
		log.Fatal(err)
	}

	// fetch access token
  // code from https://yourdomain/usercallback?code=xxxxx&shop_id=123456
  res, err := shopeeClient.Auth.GetAccessToken(shopId, 0, code)
  token := res.AccessToken

  // fetch model list
  shopeeClient.Product.GetModelList(shopId, token, 1234)

```

### Tokopedia

```
  accessToken, err := client.Auth.GetToken(ctx, clientID, secretKey)
	if err != nil {
		log.Fatal(err)
	}

	// Send a message to tokopedia
	client.Chat.SendMessage(ctx, msgId, tokopedia.TokopediaMessageText{
		Message: "This is a test message",
	})
```

## Thanks to

- [go-shopify](https://github.com/bold-commerce/go-shopify) Inspire me and provide a base structure
