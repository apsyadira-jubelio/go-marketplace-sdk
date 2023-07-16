package tokopedia

import (
	"context"
	"errors"
	"fmt"
)

type ChatService interface {
	SendMessage(ctx context.Context, msgID int64, payload TokopediaMessageText) (res *TokopediaMessageResponse, err error)
}

type ChatServiceOp struct {
	client *TokopediaClient
}

// SendMessage sends a message to the Tokopedia API.
// It accepts three parameters: a context (for managing the lifecycle of the request),
// a message ID (to identify the message), and a payload (the actual message).
// The function returns a pointer to a TokopediaMessageResponse and an error.
func (h *ChatServiceOp) SendMessage(ctx context.Context, msgID int64, payload TokopediaMessageText) (res *TokopediaMessageResponse, err error) {

	// Create the URL for the message, using the provided fsID and msgID.
	URL := fmt.Sprintf(`/v1/chat/fs/%v/messages/%v/reply`, h.client.fsID, msgID)

	// Set up and execute the request, including setting the result type, content type, and endpoint.
	resp, err := h.client.Client.R().
		SetResult(TokopediaMessageResponse{}).
		SetHeader("Content-Type", "application/json").
		Post(URL)

	if err != nil {
		return res, fmt.Errorf("%s", "Oops! gagal integrasi ke tokopedia.")
	}

	if resp.StatusCode() > 399 {
		err = errors.New(resp.String())
		return res, err
	}

	// Cast the result of the request to a TokopediaMessageResponse and assign it to the return value.
	result := resp.Result().(*TokopediaMessageResponse)
	res = result

	return res, nil
}
