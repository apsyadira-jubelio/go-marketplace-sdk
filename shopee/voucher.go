package shopee

type GetVoucherListResponse struct {
	BaseResponse
	Response VoucherList
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
