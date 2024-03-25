package tiktok

import (
	"fmt"
	"net/url"
)

type AuthService interface {
	GetAuthURL(serviceID string) (string, error)
	GetLegacyAuthURL(appKey, state string) (string, error)
	GetAccessToken(params GetAccessTokenParams) (*GetAccessTokenResponse, error)
	GetAuthorizationShop(accessToken string, shopID string) (*GetShopsResponse, error)
}

// const (
// 	AuthBaseURL   = "https://services.tiktokshop.com/open/authorize"
// 	LegacyAuthURL = "https://auth.tiktok-shops.com/oauth/authorize"
// )

type GetAccessTokenParams struct {
	AppKey    string `url:"app_key"`
	AppSecret string `url:"app_secret"`
	Code      string `url:"auth_code"`
	GrantType string `url:"grant_type"`
}

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
	aurl := fmt.Sprintf("%s/open/authorize?service_id=%s", AuthBaseURL, serviceID)
	return aurl, nil
}

func (s *AuthServiceOp) GetLegacyAuthURL(appKey, state string) (string, error) {
	aurl := fmt.Sprintf("%s/oauth/authorize?app_key=%s&state=%s", LegacyAuthURL, appKey, state)
	return aurl, nil
}

func (s *AuthServiceOp) GetAccessToken(params GetAccessTokenParams) (*GetAccessTokenResponse, error) {
	path := "/api/v2/token/get"

	resp := new(GetAccessTokenResponse)
	authURL, _ := url.Parse(LegacyAuthURL)
	s.client.baseURL = authURL
	err := s.client.Get(path, resp, params)
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

func (s *AuthServiceOp) GetAuthorizationShop(accessToken string, shopID string) (*GetShopsResponse, error) {
	// host https://open-api.tiktokglobalshop.com, automatically add app_key, sign, and timestamp in query param. Check func makeSignature
	path := fmt.Sprintf("/authorization/%s/shops", s.client.appConfig.Version)
	resp := new(GetShopsResponse)
	err := s.client.WithShopID(shopID).WithAccessToken(accessToken).Get(path, resp, nil)
	return resp, err
}
