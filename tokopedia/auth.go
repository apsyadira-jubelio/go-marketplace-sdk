package tokopedia

import (
	"context"
	"errors"
	"fmt"
)

type AuthService interface {
	GetToken(ctx context.Context, clientID string, secret string) (res *TokopediaAuthResponse, err error)
}

type AuthServiceOp struct {
	client *TokopediaClient
}

// Auth performs an authentication operation with the Tokopedia API.
// It accepts two parameters: a context (for managing the lifecycle of the request), and data (which contains authentication credentials).
// The function returns a pointer to a TokopediaAuthResponse and an error.
func (h *AuthServiceOp) GetToken(ctx context.Context, clientID string, secret string) (res *TokopediaAuthResponse, err error) {

	// Base64 encode the client ID and user secret for basic authentication
	token := Base64Encode(clientID + ":" + secret)

	// Set up and execute the request, including setting the result type, content type,
	// authorization header, and endpoint.
	resp, err := h.client.Client.R().
		SetResult(TokopediaAuthResponse{}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Basic "+token).
		Post("/token?grant_type=client_credentials")

	if err != nil {
		return res, fmt.Errorf("%s", "Oops! gagal integrasi ke tokopedia.")
	}

	// If the response has a status code of 400 or above, create and return an error.
	if resp.StatusCode() > 399 {
		err = errors.New(resp.String())
		return res, err
	}

	// Cast the result of the request to a TokopediaAuthResponse and assign it to the return value.
	result := resp.Result().(*TokopediaAuthResponse)
	res = result

	return res, nil
}
