package tokopedia

import (
	"github.com/go-resty/resty/v2"
)

func getClient(URL, APIKEY string) *resty.Client {
	client := resty.New().
		SetBaseURL(URL)

	if APIKEY != "" {
		client = client.SetHeader("Authorization", "Bearer "+APIKEY)
	}

	return client
}

func HTTPTokopediaApi(isAuth bool) *resty.Client {
	baseURL := ""

	if isAuth {
		baseURL = "https://accounts.tokopedia.com"
	} else {
		baseURL = "https://fs.tokopedia.net/v1"
	}

	return getClient(baseURL, "")
}
