package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tokopedia"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_GetToken(t *testing.T) {
	setup()
	defer teardown()

	clientID := "1828201291"
	secret := "asjdjsadiwuqeuiurekwjrwe"

	// Encode client ID and secret in base64
	token := tokopedia.Base64Encode(clientID + ":" + secret)

	// Mocking the POST request
	httpmock.RegisterResponder(http.MethodPost, "/token?grant_type=client_credentials",
		func(req *http.Request) (*http.Response, error) {
			// Create a new response
			resp := httpmock.NewBytesResponse(200, loadFixture("access_token.json"))
			resp.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Basic "+token)
			return resp, nil
		})

	// Initiate AuthServiceOp with httpmock as the Client
	client := tokopedia.NewTokopediaApi(true, "", nil)

	// Call GetToken function
	_, err := client.Auth.GetToken(context.Background(), clientID, secret)

	// Assertions
	assert.NoError(t, err)
}
