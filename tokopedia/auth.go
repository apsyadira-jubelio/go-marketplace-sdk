package tokopedia

import (
	"net/url"
)

type TokopediaAuthResponse struct {
	AccessToken   string `json:"access_token"`
	EventCode     string `json:"event_code"`
	ExpiresIn     int    `json:"expires_in"`
	LastLoginType string `json:"last_login_type"`
	SqCheck       bool   `json:"sq_check"`
	TokenType     string `json:"token_type"`
}

type AuthService interface {
	GetToken(clientID string, secret string) (res *TokopediaAuthResponse, err error)
}

type AuthServiceOp struct {
	client *TokopediaClient
}

// Auth performs an authentication operation with the Tokopedia API.
// It accepts two parameters: a context (for managing the lifecycle of the request), and data (which contains authentication credentials).
// The function returns a pointer to a TokopediaAuthResponse and an error.
func (s *AuthServiceOp) GetToken(clientID, secret string) (*TokopediaAuthResponse, error) {

	// Encode client ID and secret in base64
	var token string
	if clientID != "" && secret != "" {
		token = Base64Encode(clientID + ":" + secret)
	} else {
		token = Base64Encode(s.client.appConfig.ClientID + ":" + s.client.appConfig.ClientSecret)
	}

	path := "/token?grant_type=client_credentials"
	resp := new(TokopediaAuthResponse)

	authURL, err := url.Parse(AuthURL)
	if err != nil {
		return nil, err
	}

	s.client.baseURL = authURL
	err = s.client.WithBasicAuth(token).Post(path, nil, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
