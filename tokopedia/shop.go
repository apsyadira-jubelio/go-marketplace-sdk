package tokopedia

import "fmt"

type ShopResponse struct {
	BaseResponse
	Data []ShopData `json:"data"`
}

type ShopData struct {
	ShopID          int    `json:"shop_id"`
	UserID          int    `json:"user_id"`
	ShopName        string `json:"shop_name"`
	Logo            string `json:"logo"`
	ShopURL         string `json:"shop_url"`
	IsOpen          int    `json:"is_open"`
	Status          int    `json:"status"`
	DateShopCreated string `json:"date_shop_created"`
	Domain          string `json:"domain"`
	AdminID         []int  `json:"admin_id"`
	Reason          string `json:"reason"`
	DistrictID      int    `json:"district_id"`
	ProvinceName    string `json:"province_name"`
	Warehouses      []struct {
		WarehouseID int `json:"warehouse_id"`
		PartnerID   struct {
			Int64 int  `json:"Int64"`
			Valid bool `json:"Valid"`
		} `json:"partner_id"`
		ShopID struct {
			Int64 int  `json:"Int64"`
			Valid bool `json:"Valid"`
		} `json:"shop_id"`
		WarehouseName string `json:"warehouse_name"`
		DistrictID    int    `json:"district_id"`
		DistrictName  string `json:"district_name"`
		CityID        int    `json:"city_id"`
		CityName      string `json:"city_name"`
		ProvinceID    int    `json:"province_id"`
		ProvinceName  string `json:"province_name"`
		Status        int    `json:"status"`
		PostalCode    string `json:"postal_code"`
		IsDefault     int    `json:"is_default"`
		Latlon        string `json:"latlon"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
		Email         string `json:"email"`
		AddressDetail string `json:"address_detail"`
		Phone         string `json:"phone"`
		WarehoseType  string `json:"warehose_type"`
	} `json:"warehouses"`
	SubscribeTokocabang bool `json:"subscribe_tokocabang"`
	IsMitra             bool `json:"is_mitra"`
}

type ShopParams struct {
	ShopID  int `url:"shop_id"`
	Page    int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`
}

type ShopService interface {
	GetShopInfo(token string, params ShopParams) (res *ShopResponse, err error)
}

type ShopServiceOp struct {
	client *TokopediaClient
}

func (s *ShopServiceOp) GetShopInfo(token string, params ShopParams) (res *ShopResponse, err error) {
	path := fmt.Sprintf("/v1/shop/fs/%d/shop-info", s.client.appConfig.FsID)
	resp := new(ShopResponse)

	err = s.client.WithAccessToken(token).Get(path, resp, params)
	return resp, err
}
