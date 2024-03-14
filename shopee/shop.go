package shopee

type ShopService interface {
	GetShopInfo(shopID uint64, token string) (*GetShopInfoResponse, error)
}

type GetShopInfoResponse struct {
	ShopName            string `json:"shop_name"`
	Region              string `json:"region"`
	Status              string `json:"status"`
	IsSip               bool   `json:"is_sip"`
	IsCb                bool   `json:"is_cb"`
	IsCnsc              bool   `json:"is_cnsc"`
	RequestID           string `json:"request_id"`
	AuthTime            int    `json:"auth_time"`
	ExpireTime          int    `json:"expire_time"`
	Error               string `json:"error"`
	Message             string `json:"message"`
	ShopCbsc            string `json:"shop_cbsc"`
	MtskuUpgradedStatus string `json:"mtsku_upgraded_status"`
}

type ShopServiceOp struct {
	client *ShopeeClient
}

func (s *ShopServiceOp) GetShopInfo(shopID uint64, token string) (*GetShopInfoResponse, error) {
	path := "/shop/get_shop_info"
	resp := new(GetShopInfoResponse)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, nil)
	return resp, err
}
