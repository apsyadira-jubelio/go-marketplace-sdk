package tests

import (
	"context"
	"testing"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/lazada"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_SendMessage(t *testing.T) {
	setup()
	defer teardown()

	// sessionID := "400096074640_1_400601424036_2_103"
	// templateID := 1
	// txt := "Hello"
	// ApiURL := fmt.Sprintf("%s?session_id=%s&template_id=%d&txt=%s", lazada.ApiNames["SendMessage"], sessionID, templateID, txt)

	httpmock.RegisterResponder("POST",
		`^.*\/im\/message\/send.*$`,
		httpmock.NewBytesResponder(200, loadFixture("send_message_resp.json")),
	)

	var req lazada.SendMessageParams
	loadMockData("send_message_req.json", &req)

	res, err := client.Chat.SendMessage(context.TODO(), "50000600116cWYzTphDtTDshMBux1993574eoq9YzkugHtfWTiXeDQ7OzvDLRkFx", &lazada.SendMessageParams{
		SessionID:  req.SessionID,
		TemplateID: req.TemplateID,
		Txt:        req.Txt,
	})
	// Check that there was no error
	assert.Nil(t, err)

	if res != nil {
		assert.Equal(t, true, res.Success)
	}
}
