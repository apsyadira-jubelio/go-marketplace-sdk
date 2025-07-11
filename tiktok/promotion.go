package tiktok

type PromotionService interface {
	SearchCoupons(body SearchCouponsBody) (*SearchCouponsResponse, error)
}

type PromotionServiceOp struct {
	client *TiktokClient
}

type SearchCouponsBody struct {
	Status       []string `json:"status"`
	TitleKeyword string   `json:"title_keyword,omitempty"`
	DisplayType  []string `json:"display_type,omitempty"`
}

type SearchCouponsResponse struct {
	BaseResponse
	Data DataCoupon `json:"data"`
}

type ClaimDuration struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
}

type RedemptionDuration struct {
	Type         string `json:"type"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	RelativeTime int    `json:"relative_time"`
}

type UsageLimits struct {
	SingleBuyerClaimLimit int `json:"single_buyer_claim_limit"`
	TotalClaimLimit       int `json:"total_claim_limit"`
	RedemptionLimit       int `json:"redemption_limit"`
}

type ReductionAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type MaxDiscount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Discount struct {
	Type            string          `json:"type"`
	ReductionAmount ReductionAmount `json:"reduction_amount"`
	Percentage      string          `json:"percentage"`
	MaxDiscount     MaxDiscount     `json:"max_discount"`
}

type MinSpend struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Threshold struct {
	Type     string   `json:"type"`
	MinSpend MinSpend `json:"min_spend"`
}

type Coupons struct {
	ID                 string             `json:"id"`
	Title              string             `json:"title"`
	DisplayType        string             `json:"display_type"`
	Status             string             `json:"status"`
	CreateTime         int64              `json:"create_time"`
	UpdateTime         int64              `json:"update_time"`
	ClaimDuration      ClaimDuration      `json:"claim_duration"`
	RedemptionDuration RedemptionDuration `json:"redemption_duration"`
	PromoCode          string             `json:"promo_code"`
	TargetBuyerSegment string             `json:"target_buyer_segment"`
	UsageLimits        UsageLimits        `json:"usage_limits"`
	Discount           Discount           `json:"discount"`
	Threshold          Threshold          `json:"threshold"`
	ProductScope       string             `json:"product_scope"`
	CreationSource     string             `json:"creation_source"`
}
type DataCoupon struct {
	TotalCount    int       `json:"total_count"`
	NextPageToken string    `json:"next_page_token"`
	Coupons       []Coupons `json:"coupons"`
}

func (s *PromotionServiceOp) SearchCoupons(body SearchCouponsBody) (*SearchCouponsResponse, error) {
	path := "/promotion/202406/coupons/search?page_size=100"
	resp := new(SearchCouponsResponse)
	err := s.client.Post(path, body, resp)
	return resp, err
}
