package lazada

import (
	"context"
	"encoding/json"
)

// The Product Service deals with any methods under the "Instant Messaging" category of the open platform
type ProductService service

type GetProductsParams struct {
	Filter       string `url:"filter,omitempty"`
	Offset       string `url:"offset"`
	Limit        string `url:"limit"`
	CreateBefore string `url:"create_before"`
	UpdateBefore string `url:"update_before"`
	CreateAfter  string `url:"create_after"`
	UpdateAfter  string `url:"update_after"`
}

type GetProductsResponse struct {
	BaseResponse
	Data DataGetProducts `json:"data"`
}

type DataGetProducts struct {
	TotalProducts int        `json:"total_products"`
	Products      []Products `json:"products"`
}

type Products struct {
	CreatedTime     interface{}    `json:"created_time"`
	UpdatedTime     interface{}    `json:"updated_time"`
	Images          []string       `json:"images"`
	Skus            []Skus         `json:"skus"`
	ItemID          int            `json:"item_id"`
	HiddenStatus    string         `json:"hiddenStatus"`
	SuspendedSkus   []interface{}  `json:"suspendedSkus"`
	SubStatus       string         `json:"subStatus"`
	TrialProduct    string         `json:"trialProduct"`
	RejectReason    []RejectReason `json:"rejectReason"`
	PrimaryCategory string         `json:"primary_category"`
	MarketImages    string         `json:"marketImages"`
	Attributes      Attributes     `json:"attributes"`
	HiddenReason    string         `json:"hiddenReason"`
	Status          string         `json:"status"`
}

type Skus struct {
	Status          string   `json:"Status"`
	Quantity        int      `json:"quantity"`
	ProductWeight   int      `json:"product_weight"`
	Images          []string `json:"Images"`
	SellerSku       string   `json:"SellerSku"`
	ShopSku         string   `json:"ShopSku"`
	URL             string   `json:"Url"`
	PackageWidth    string   `json:"package_width"`
	SpecialToTime   string   `json:"special_to_time"`
	SpecialFromTime string   `json:"special_from_time"`
	PackageHeight   string   `json:"package_height"`
	SpecialPrice    float64  `json:"special_price"`
	Price           float64  `json:"price"`
	PackageLength   string   `json:"package_length"`
	PackageWeight   string   `json:"package_weight"`
	Available       int      `json:"Available"`
	SkuID           int      `json:"SkuId"`
	SpecialToDate   string   `json:"special_to_date"`
}

type RejectReason struct {
	Suggestion      string `json:"suggestion"`
	ViolationDetail string `json:"violationDetail"`
}

type Attributes struct {
	ShortDescription string `json:"short_description"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	NameEngravement  string `json:"name_engravement"`
	WarrantyType     string `json:"warranty_type"`
	GiftWrapping     string `json:"gift_wrapping"`
	PreorderDays     int    `json:"preorder_days"`
	Brand            string `json:"brand"`
	Preorder         string `json:"preorder"`
}

// GetProducts is a method on the ProductService struct. Use this API to get detailed information of the specified products.
// If the opts parameter is nil, default options will be used with a 50 limit products with no filter.
// The function returns a pointer to an GetProductsResponse struct containing the server's response, and an error, if there is one.
func (p *ProductService) GetProducts(ctx context.Context, opts *GetProductsParams) (res *GetProductsResponse, err error) {
	if opts == nil {
		opts = &GetProductsParams{
			Limit:  "25",
			Offset: "0",
		}
	}

	u, err := addOptions(ApiNames["GetProducts"], opts)
	if err != nil {
		return nil, err
	}

	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(jsonData), &res)

	return res, nil
}
