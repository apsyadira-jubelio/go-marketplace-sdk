package shopee

type LogisticService interface {
	GetTrackingInfo(shopID uint64, token string, params GetTrackingInfoParamsRequest) (*GetTrackingInfoResponse, error)
}

type GetTrackingInfoParamsRequest struct {
	OrderSN       string `url:"order_sn"` // based on doc, it's string
	PackageNumber string `url:"package_number,omitempty"`
}

// Response Get Detail Order
type (
	GetTrackingInfoResponse struct {
		BaseResponse
		GetTrackingInfoData any `json:"response"`
	}
)

type LogisticServiceOp struct {
	client *ShopeeClient
}

func (o *LogisticServiceOp) GetTrackingInfo(shopID uint64, token string, params GetTrackingInfoParamsRequest) (*GetTrackingInfoResponse, error) {
	path := "/logistics/get_tracking_info"
	resp := new(GetTrackingInfoResponse)
	err := o.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}
