package tiktok

import (
	"fmt"
)

type OrderService interface {
	GetOrder(params GetOrderParams) (*GetOrderResponse, error)
}

type OrderServiceOp struct {
	client *TiktokClient
}

type GetOrderParams struct {
	OrderIDs []string `url:"ids"`
}

type GetOrderResponse struct {
	BaseResponse
	Data OrderData `json:"data"`
}

type OrderData struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	BuyerEmail                         string           `json:"buyer_email"`
	BuyerMessage                       string           `json:"buyer_message"`
	CancelOrderSlaTime                 int64            `json:"cancel_order_sla_time"`
	CancelReason                       string           `json:"cancel_reason"`
	CancelTime                         int64            `json:"cancel_time"`
	CancellationInitiator              string           `json:"cancellation_initiator"`
	CollectionDueTime                  int64            `json:"collection_due_time"`
	CollectionTime                     int64            `json:"collection_time"`
	Cpf                                string           `json:"cpf"`
	CreateTime                         int64            `json:"create_time"`
	DeliveryDueTime                    int64            `json:"delivery_due_time"`
	DeliveryOptionID                   string           `json:"delivery_option_id"`
	DeliveryOptionName                 string           `json:"delivery_option_name"`
	DeliveryOptionRequiredDeliveryTime int64            `json:"delivery_option_required_delivery_time"`
	DeliverySlaTime                    int64            `json:"delivery_sla_time"`
	DeliveryTime                       int64            `json:"delivery_time"`
	FulfillmentType                    string           `json:"fulfillment_type"`
	HasUpdatedRecipientAddress         bool             `json:"has_updated_recipient_address"`
	ID                                 string           `json:"id"`
	IsBuyerRequestCancel               bool             `json:"is_buyer_request_cancel"`
	IsCod                              bool             `json:"is_cod"`
	IsOnHoldOrder                      bool             `json:"is_on_hold_order"`
	IsReplacementOrder                 bool             `json:"is_replacement_order"`
	IsSampleOrder                      bool             `json:"is_sample_order"`
	LineItems                          []LineItem       `json:"line_items"`
	NeedUploadInvoice                  string           `json:"need_upload_invoice"`
	Packages                           []Package        `json:"packages"`
	PaidTime                           int64            `json:"paid_time"`
	Payment                            Payment          `json:"payment"`
	PaymentMethodName                  string           `json:"payment_method_name"`
	RecipientAddress                   RecipientAddress `json:"recipient_address"`
	ReplacedOrderID                    string           `json:"replaced_order_id"`
	RequestCancelTime                  int64            `json:"request_cancel_time"`
	RTSSlaTime                         int64            `json:"rts_sla_time"`
	RTSTime                            int64            `json:"rts_time"`
	SellerNote                         string           `json:"seller_note"`
	ShippingDueTime                    int64            `json:"shipping_due_time"`
	ShippingProvider                   string           `json:"shipping_provider"`
	ShippingProviderID                 string           `json:"shipping_provider_id"`
	ShippingType                       string           `json:"shipping_type"`
	SplitOrCombineTag                  string           `json:"split_or_combine_tag"`
	Status                             string           `json:"status"`
	TrackingNumber                     string           `json:"tracking_number"`
	TTSSlaTime                         int64            `json:"tts_sla_time"`
	UpdateTime                         int64            `json:"update_time"`
	UserID                             string           `json:"user_id"`
	WarehouseID                        string           `json:"warehouse_id"`
}

type LineItem struct {
	CancelReason         string                `json:"cancel_reason"`
	CancelUser           string                `json:"cancel_user"`
	CombinedListingSkus  []CombinedListingSkus `json:"combined_listing_skus"`
	Currency             string                `json:"currency"`
	DisplayStatus        string                `json:"display_status"`
	ID                   string                `json:"id"`
	IsGift               bool                  `json:"is_gift"`
	ItemTax              []ItemTax             `json:"item_tax"`
	OriginalPrice        string                `json:"original_price"`
	PackageID            string                `json:"package_id"`
	PackageStatus        string                `json:"package_status"`
	PlatformDiscount     string                `json:"platform_discount"`
	ProductID            string                `json:"product_id"`
	ProductName          string                `json:"product_name"`
	RetailDeliveryFee    string                `json:"retail_delivery_fee"`
	RTSTime              int64                 `json:"rts_time"`
	SalePrice            string                `json:"sale_price"`
	SellerDiscount       string                `json:"seller_discount"`
	SellerSku            string                `json:"seller_sku"`
	ShippingProviderID   string                `json:"shipping_provider_id"`
	ShippingProviderName string                `json:"shipping_provider_name"`
	SkuID                string                `json:"sku_id"`
	SkuImage             string                `json:"sku_image"`
	SkuName              string                `json:"sku_name"`
	SkuType              string                `json:"sku_type"`
	SmallOrderFee        string                `json:"small_order_fee"`
	TrackingNumber       string                `json:"tracking_number"`
}

type CombinedListingSkus struct {
	ProductID string `json:"product_id"`
	SkuCount  int64  `json:"sku_count"`
	SkuID     string `json:"sku_id"`
}

type ItemTax struct {
	TaxAmount string `json:"tax_amount"`
	TaxRate   string `json:"tax_rate"`
	TaxType   string `json:"tax_type"`
}

type Package struct {
	ID string `json:"id"`
}

type Payment struct {
	Currency                    string `json:"currency"`
	OriginalShippingFee         string `json:"original_shipping_fee"`
	OriginalTotalProductPrice   string `json:"original_total_product_price"`
	PlatformDiscount            string `json:"platform_discount"`
	ProductTax                  string `json:"product_tax"`
	RetailDeliveryFee           string `json:"retail_delivery_fee"`
	SellerDiscount              string `json:"seller_discount"`
	ShippingFee                 string `json:"shipping_fee"`
	ShippingFeePlatformDiscount string `json:"shipping_fee_platform_discount"`
	ShippingFeeSellerDiscount   string `json:"shipping_fee_seller_discount"`
	ShippingFeeTax              string `json:"shipping_fee_tax"`
	SmallOrderFee               string `json:"small_order_fee"`
	SubTotal                    string `json:"sub_total"`
	Tax                         string `json:"tax"`
	TotalAmount                 string `json:"total_amount"`
}

type RecipientAddress struct {
	AddressDetail       string              `json:"address_detail"`
	AddressLine1        string              `json:"address_line1"`
	AddressLine2        string              `json:"address_line2"`
	AddressLine3        string              `json:"address_line3"`
	AddressLine4        string              `json:"address_line4"`
	DeliveryPreferences DeliveryPreferences `json:"delivery_preferences"`
	DistrictInfo        []DistrictInfo      `json:"district_info"`
	FullAddress         string              `json:"full_address"`
	Name                string              `json:"name"`
	PhoneNumber         string              `json:"phone_number"`
	PostalCode          string              `json:"postal_code"`
	RegionCode          string              `json:"region_code"`
}

type DeliveryPreferences struct {
	DropOffLocation string `json:"drop_off_location"`
}

type DistrictInfo struct {
	AddressLevel     string `json:"address_level"`
	AddressLevelName string `json:"address_level_name"`
	AddressName      string `json:"address_name"`
}

func (s *OrderServiceOp) GetOrder(params GetOrderParams) (*GetOrderResponse, error) {
	path := fmt.Sprintf("/order/%s/orders", s.client.appConfig.Version)

	resp := new(GetOrderResponse)
	err := s.client.Get(path, resp, params)

	return resp, err
}
