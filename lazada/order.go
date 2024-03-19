package lazada

import (
	"context"
	"encoding/json"
)

// The Order Service deals with any methods under the "Order" category of the open platform
type OrderService service

type GetMultipleOrdersItemsParam struct {
	OrderIDs []int `url:"order_ids"`
}

type GetMultipleOrdersItemsResponse struct {
	BaseResponse
	DataOrder []DataGetMultipleOrderItem `json:"data"`
}

type DataGetMultipleOrderItem struct {
	OrderID     int64        `json:"order_id"`
	OrderItems  []OrderItems `json:"order_items"`
	OrderNumber int64        `json:"order_number"`
}

type PickUpStoreInfo struct {
}

type OrderItems struct {
	BuyerID                     int64           `json:"buyer_id"`
	CancelReturnInitiator       string          `json:"cancel_return_initiator"`
	CreatedAt                   string          `json:"created_at"`
	Currency                    string          `json:"currency"`
	DeliveryOptionSof           int             `json:"delivery_option_sof"`
	DigitalDeliveryInfo         string          `json:"digital_delivery_info"`
	ExtraAttributes             string          `json:"extra_attributes"`
	FulfillmentSLA              string          `json:"fulfillment_sla"`
	GiftWrapping                string          `json:"gift_wrapping"`
	InvoiceNumber               string          `json:"invoice_number"`
	IsDigital                   int             `json:"is_digital"`
	IsFbl                       int             `json:"is_fbl"`
	IsReroute                   int             `json:"is_reroute"`
	ItemPrice                   int             `json:"item_price"`
	Mp3Order                    bool            `json:"mp3_order"`
	Name                        string          `json:"name"`
	OrderFlag                   string          `json:"order_flag"`
	OrderID                     int64           `json:"order_id"`
	OrderItemID                 int64           `json:"order_item_id"`
	OrderType                   string          `json:"order_type"`
	PackageID                   string          `json:"package_id"`
	PaidPrice                   int             `json:"paid_price"`
	Personalization             string          `json:"personalization"`
	PickUpStoreInfo             PickUpStoreInfo `json:"pick_up_store_info"`
	PriorityFulfillmentTag      string          `json:"priority_fulfillment_tag"`
	ProductDetailURL            string          `json:"product_detail_url"`
	ProductMainImage            string          `json:"product_main_image"`
	PromisedShippingTime        string          `json:"promised_shipping_time"`
	PurchaseOrderID             string          `json:"purchase_order_id"`
	PurchaseOrderNumber         string          `json:"purchase_order_number"`
	Reason                      string          `json:"reason"`
	ReasonDetail                string          `json:"reason_detail"`
	ReturnStatus                string          `json:"return_status"`
	SemiManaged                 bool            `json:"semi_managed"`
	ShipmentProvider            string          `json:"shipment_provider"`
	ShippingAmount              int             `json:"shipping_amount"`
	ShippingFeeDiscountPlatform int             `json:"shipping_fee_discount_platform"`
	ShippingFeeDiscountSeller   int             `json:"shipping_fee_discount_seller"`
	ShippingFeeOriginal         int             `json:"shipping_fee_original"`
	ShippingProviderType        string          `json:"shipping_provider_type"`
	ShippingServiceCost         int             `json:"shipping_service_cost"`
	ShippingType                string          `json:"shipping_type"`
	ShopID                      string          `json:"shop_id"`
	ShopSku                     string          `json:"shop_sku"`
	Sku                         string          `json:"sku"`
	SkuID                       string          `json:"sku_id"`
	SLATimeStamp                string          `json:"sla_time_stamp"`
	StagePayStatus              string          `json:"stage_pay_status"`
	Status                      string          `json:"status"`
	SupplyPrice                 int             `json:"supply_price"`
	SupplyPriceCurrency         string          `json:"supply_price_currency"`
	TaxAmount                   int             `json:"tax_amount"`
	TrackingCode                string          `json:"tracking_code"`
	TrackingCodePre             string          `json:"tracking_code_pre"`
	UpdatedAt                   string          `json:"updated_at"`
	Variation                   string          `json:"variation"`
	VoucherAmount               int             `json:"voucher_amount"`
	VoucherCode                 string          `json:"voucher_code"`
	VoucherCodePlatform         string          `json:"voucher_code_platform"`
	VoucherCodeSeller           string          `json:"voucher_code_seller"`
	VoucherPlatform             int             `json:"voucher_platform"`
	VoucherPlatformLpi          int             `json:"voucher_platform_lpi"`
	VoucherSeller               int             `json:"voucher_seller"`
	VoucherSellerLpi            int             `json:"voucher_seller_lpi"`
	WalletCredits               int             `json:"wallet_credits"`
	WarehouseCode               string          `json:"warehouse_code"`
}

// GetMultipleOrdersItems is a method on the OrderService struct. Use this API to get detailed information of the specified orders.
// The function returns a pointer to an GetMultipleOrdersItemsResponse struct containing the server's response, and an error, if there is one.
func (o *OrderService) GetMultipleOrdersItems(ctx context.Context, opts *GetMultipleOrdersItemsParam) (res *GetMultipleOrdersItemsResponse, err error) {
	u, err := addOptions(ApiNames["GetMultipleOrdersItems"], opts)
	if err != nil {
		return nil, err
	}

	req, err := o.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := o.client.Do(ctx, req, nil)
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
