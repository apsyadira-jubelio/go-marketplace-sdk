package shopee

type VoucherService interface {
	GetListVoucherByStatus(shopID uint64, token string, params GetVoucherListParam) (*GetVoucherListResponse, error)
	GetDetailVoucher(shopID uint64, token string, params GetDetailVoucherParam) (any, error)
}

type GetVoucherListParam struct {
	PageNo   int    `url:"page_no"`
	PageSize int    `url:"page_size"`
	Status   string `url:"status"` // required. Available value: upcoming/ongoing/expired/all.
}

type GetDetailVoucherParam struct {
	VoucherID int64 `url:"voucher_id"`
}

type GetVoucherListResponse struct {
	BaseResponse
	Response VoucherListResponse
}

type VoucherListResponse struct {
	More        bool          `json:"more"`
	VoucherList []VoucherList `json:"voucher_list"`
}

type VoucherList struct {
	VoucherID        int64  `json:"voucher_id"`
	VoucherCode      string `json:"voucher_code"`
	VoucherName      string `json:"voucher_name"`
	VoucherType      int    `json:"voucher_type"`
	RewardType       int    `json:"reward_type"`
	UsageQuantity    int    `json:"usage_quantity"`
	CurrentUsage     int    `json:"current_usage"`
	StartTime        int    `json:"start_time"`
	EndTime          int    `json:"end_time"`
	IsAdmin          bool   `json:"is_admin"`
	VoucherPurpose   int    `json:"voucher_purpose"`
	DiscountAmount   int    `json:"discount_amount,omitempty"`
	TargetVoucher    int    `json:"target_voucher"`
	DisplayStartTime int    `json:"display_start_time"`
	Percentage       int    `json:"percentage,omitempty"`
}

type VoucherServiceOp struct {
	client *ShopeeClient
}

func (v *VoucherServiceOp) GetListVoucherByStatus(shopID uint64, token string, params GetVoucherListParam) (*GetVoucherListResponse, error) {
	path := "/voucher/get_voucher_list"
	resp := new(GetVoucherListResponse)
	err := v.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}

func (v *VoucherServiceOp) GetDetailVoucher(shopID uint64, token string, params GetDetailVoucherParam) (any, error) {
	path := "/voucher/get_voucher"
	resp := new(any)
	err := v.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}
