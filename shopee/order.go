package shopee

type OrderService interface {
	GetOrderDetailByOrderSN(shopID uint64, token string, params GetOrderDetailParamsRequest) (*GetOrderDetailResponse, error)
	GetListOrder(shopID uint64, token string, params GetListOrderParamsRequest) (*GetListOrderResponse, error)
	DownloadInvoiceByOrderID(shopID uint64, token string, params DownloadInvoiceParamsRequest) error
}

type GetOrderDetailParamsRequest struct {
	OrderSNList               string `url:"order_sn_list"` // based on doc, it's string
	RequestOrderStatusPending bool   `url:"request_order_status_pending,omitempty"`
	ResponseOptionalFields    string `url:"response_optional_fields"`
}

// Response Get Detail Order
type (
	GetOrderDetailResponse struct {
		BaseResponse
		OrderListResponse OrderListResponse `json:"response"`
	}

	OrderListResponse struct {
		OrderList []OrderList `json:"order_list"`
	}

	ItemList struct {
		ItemID                 int64      `json:"item_id"`
		ItemName               string     `json:"item_name"`
		ItemSku                string     `json:"item_sku"`
		ModelID                int64      `json:"model_id"`
		ModelName              string     `json:"model_name"`
		ModelSku               string     `json:"model_sku"`
		ModelQuantityPurchased int        `json:"model_quantity_purchased"`
		ModelOriginalPrice     int        `json:"model_original_price"`
		ModelDiscountedPrice   int        `json:"model_discounted_price"`
		Wholesale              bool       `json:"wholesale"`
		Weight                 float64    `json:"weight"`
		AddOnDeal              bool       `json:"add_on_deal"`
		MainItem               bool       `json:"main_item"`
		AddOnDealID            int        `json:"add_on_deal_id"`
		PromotionType          string     `json:"promotion_type"`
		PromotionID            int        `json:"promotion_id"`
		OrderItemID            int64      `json:"order_item_id"`
		PromotionGroupID       int        `json:"promotion_group_id"`
		ImageInfo              *ImageInfo `json:"image_info"`
		ProductLocationID      []string   `json:"product_location_id"`
		IsPrescriptionItem     bool       `json:"is_prescription_item"`
		IsB2COwnedItem         bool       `json:"is_b2c_owned_item"`
	}

	OrderList struct {
		Cod                bool        `json:"cod"`
		CreateTime         int         `json:"create_time"`
		Currency           string      `json:"currency"`
		DaysToShip         int         `json:"days_to_ship"`
		ItemList           []ItemList  `json:"item_list"`
		InvoiceData        interface{} `json:"invoice_info_list"`
		MessageToSeller    string      `json:"message_to_seller"`
		OrderSn            string      `json:"order_sn"`
		OrderStatus        string      `json:"order_status"`
		PaymentMethod      string      `json:"payment_method"`
		Region             string      `json:"region"`
		ReverseShippingFee int         `json:"reverse_shipping_fee"`
		ShipByDate         int         `json:"ship_by_date"`
		ShippingCarrier    string      `json:"shipping_carrier"`
		TotalAmount        int         `json:"total_amount"`
		UpdateTime         int         `json:"update_time"`
	}
)

type OrderServiceOp struct {
	client *ShopeeClient
}

func (o *OrderServiceOp) GetOrderDetailByOrderSN(shopID uint64, token string, params GetOrderDetailParamsRequest) (*GetOrderDetailResponse, error) {
	path := "/order/get_order_detail"
	params.ResponseOptionalFields = "buyer_user_id,buyer_username,estimated_shipping_fee,recipient_address,actual_shipping_fee,goods_to_declare,note,note_update_time,item_list,pay_time,dropshipper, dropshipper_phone,split_up,buyer_cancel_reason,cancel_by,cancel_reason,actual_shipping_fee_confirmed,buyer_cpf_id,fulfillment_flag,pickup_done_time,package_list,shipping_carrier,payment_method,total_amount,buyer_username,invoice_info_list,no_plastic_packing,order_chargeable_weight_gram,edt,return_due_date"
	resp := new(GetOrderDetailResponse)
	err := o.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}

type DownloadInvoiceParamsRequest struct {
	OrderSNList string `url:"order_sn_list"`
}

func (o *OrderServiceOp) DownloadInvoiceByOrderID(shopID uint64, token string, params DownloadInvoiceParamsRequest) error {
	path := "/order/download_invoice_doc"
	err := o.client.WithShop(uint64(shopID), token).Get(path, "", params)
	return err
}

type GetListOrderParamsRequest struct {
	TimeRangeField string `url:"time_range_field"` // create/update time
	TimeFrom       int    `url:"time_from"`        // epoch based
	TimeTo         int    `url:"time_to"`
	PageSize       int    `url:"page_size"`
	OrderStatus    string `url:"order_status,omitempty"` // UNPAID/READY_TO_SHIP/PROCESSED/SHIPPED/COMPLETED/IN_CANCEL/CANCELLED/INVOICE_PENDING
}

type GetListOrderResponse struct {
	BaseResponse
	Response ResponseOrderList `json:"response"`
}

type ResponseOrderList struct {
	More       bool          `json:"more"`
	NextCursor string        `json:"next_cursor"`
	OrderList  []OrderSNList `json:"order_list"`
}
type OrderSNList struct {
	OrderSn string `json:"order_sn"`
}

func (o *OrderServiceOp) GetListOrder(shopID uint64, token string, params GetListOrderParamsRequest) (*GetListOrderResponse, error) {
	path := "/order/get_order_list"
	resp := new(GetListOrderResponse)
	err := o.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}
