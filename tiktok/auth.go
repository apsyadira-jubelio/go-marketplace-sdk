package tiktok

import (
	"fmt"
)

type AuthService interface {
	GetAuthURL(serviceID string) (string, error)
	GetOldAuthURL(appKey, state string) (string, error)
	GetAccessToken(appKey, appSecret, code, grantType string) (*GetAccessTokenResponse, error)
	GetAuthorizationShop(version string) (*GetShopsResponse, error)
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

func (s *AuthServiceOp) GetAccessToken(appKey, appSecret, code, grantType string) (*GetAccessTokenResponse, error) {
	path := fmt.Sprintf("/api/v2/token/get?app_key%s&app_secret=%s&auth_code=%s&grant_type=%s", appKey, appSecret, code, grantType)
	resp := new(GetAccessTokenResponse)
	err := s.client.Get(path, nil, resp)
	return resp, err
}

type GetShopsResponse struct {
	BaseResponse
	Data DataShops `json:"data"`
}
type Shops struct {
	Cipher     string `json:"cipher"`
	Code       string `json:"code"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Region     string `json:"region"`
	SellerType string `json:"seller_type"`
}
type DataShops struct {
	Shops []Shops `json:"shops"`
}

func (s *AuthServiceOp) GetAuthorizationShop(version string) (*GetShopsResponse, error) {
	// host https://open-api.tiktokglobalshop.com, automatically add app_key, sign, and timestamp in query param. Check func makeSignature
	path := fmt.Sprintf("/authorization/%s/shops", version)
	resp := new(GetShopsResponse)
	err := s.client.Get(path, nil, resp)
	return resp, err
}
