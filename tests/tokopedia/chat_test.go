package tests

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tokopedia"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_SendMessage(t *testing.T) {
	setup()
	defer teardown()

	token := "c:hd5JHwOcRk2QSSAde_Z2bw"
	strVal := "12345"
	fsID, _ := strconv.ParseInt(strVal, 10, 64)

	// Mocking the POST request
	httpmock.RegisterResponder(http.MethodPost, "/chat/fs/12345/messages/54321/reply",
		func(req *http.Request) (*http.Response, error) {
			// Create a new response
			resp := httpmock.NewBytesResponse(200, loadFixture("send_message_resp.json"))
			resp.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			return resp, nil
		})

	// Initiate AuthServiceOp with httpmock as the Client
	client := tokopedia.NewClient(false, token, &fsID)

	var req tokopedia.TokopediaMessageText
	loadMockData("send_message_req.json", &req)

	// Call GetToken function
	_, err := client.Chat.SendMessage(context.Background(), 54321, tokopedia.TokopediaMessageText{
		ShopId:  req.ShopId,
		Message: req.Message,
	})

	// Assertions
	assert.NoError(t, err)
}
