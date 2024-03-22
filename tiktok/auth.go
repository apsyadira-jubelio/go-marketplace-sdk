package tiktok

import (
	"fmt"
)

type AuthService interface {
	GetAuthURL(serviceID string) (string, error)
	GetOldAuthURL(appKey, state string) (string, error)
	GetAccessToken(code, grantType string) (*GetAccessTokenResponse, error)
}

const (
	AuthBaseURL    = "https://services.tiktokshop.com/open/authorize"
	OldAuthBaseURL = "https://auth.tiktok-shops.com/oauth/authorize"
)

type GetAccessTokenResponse struct {
	BaseResponse
	Data DataAccessToken `json:"data"`
}

type DataAccessToken struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpireIn  int    `json:"access_token_expire_in"`
	RefreshToken         string `json:"refresh_token"`
	RefreshTokenExpireIn int    `json:"refresh_token_expire_in"`
	OpenID               string `json:"open_id"`
	SellerName           string `json:"seller_name"`
	SellerBaseRegion     string `json:"seller_base_region"`
	UserType             int    `json:"user_type"`
}

type AuthServiceOp struct {
	client *TiktokClient
}

func (s *AuthServiceOp) GetAuthURL(serviceID string) (string, error) {
	aurl := fmt.Sprintf("%s?service_id=%s", AuthBaseURL, serviceID)
	return aurl, nil
}

func (s *AuthServiceOp) GetOldAuthURL(appKey, state string) (string, error) {
	aurl := fmt.Sprintf("%s?app_key=%s&state=%s", OldAuthBaseURL, appKey, state)
	return aurl, nil
}

func (s *AuthServiceOp) GetAccessToken(code, grantType string) (*GetAccessTokenResponse, error) {
	path := fmt.Sprintf("/auth/token/get?app_secret=%s&auth_code=%s&grant_type=%s", s.client.appConfig.AppSecret, code, grantType)
	resp := new(GetAccessTokenResponse)
	err := s.client.Get(path, nil, resp)
	return resp, err
}
