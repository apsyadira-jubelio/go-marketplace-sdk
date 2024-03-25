package tests

import (
	"fmt"
	"testing"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tiktok"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_GetAuthURL(t *testing.T) {
	setup()
	defer teardown()

	authURL, _ := client.Auth.GetLegacyAuthURL(app.AppKey, "teststate")
	assert.NotEqual(t, authURL, tiktok.LegacyAuthURL, "auth url should be same")
	assert.NotEmpty(t, authURL, "auth url should not be empty")
}

func Test_GetAccessToken(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/v2/token/get", tiktok.LegacyAuthURL),
		httpmock.NewBytesResponder(200, loadFixture("access_token_resp.json")))

	fmt.Println(app.APIURL)

	qs := tiktok.GetAccessTokenParams{
		AppKey:    app.AppKey,
		AppSecret: app.AppSecret,
		Code:      "testcode",
		GrantType: "authorization_code",
	}
	res, err := client.Auth.GetAccessToken(qs)
	if err != nil {
		t.Errorf("Auth.GetToken error: %s", err)
	}

	var expectedToken string = "ROW_-9D27wAAAABftY_-lBYbKUNezeTwBEzV7T-uEdQR3qD7lu7tdl0YuX1OsYoBtH2L1nlzgH-m4OYORtNg3YKqUPBdiuleV17Tnndh8v9jpeM4Zk-pinJ7V53fA7DWDVJHoD2f-2YqfN63WY6BhtuaNrlhddEr6ZVZr3osL1nQogHjLgU6Hfs4CA"
	assert.Equal(t, expectedToken, res.Data.AccessToken, "Data.AccessToken should be equal")
	assert.NotEqual(t, res.Data.AccessToken, "", "Data.AccessToken should not be empty")
	assert.NotEqualValues(t, expectedToken, "12", "Data.AccessToken should not be 12319381092")

	t.Logf("return tok: %#v", res)
}
