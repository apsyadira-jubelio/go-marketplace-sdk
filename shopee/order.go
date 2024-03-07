package shopee

type OrderService interface {
	GetOrderDetailByOrderSN(shopID uint64, token string, params GetOrderDetailParamsRequest) (*GetOrderDetailResponse, error)
}

type GetOrderDetailParamsRequest struct {
	OrderSNList               string `url:"order_sn_list"` // based on doc, it's string
	RequestOrderStatusPending bool   `url:"request_order_status_pending,omitempty"`
	ResponseOptionalFields    string `url:"response_optional_fields,omitempty"`
}

// Response Get Detail Order
type (
	GetOrderDetailResponse struct {
		BaseResponse
		InvoiceInfoList []InvoiceInfoList `json:"invoice_info_list"`
	}
	AddressBreakdown struct {
		Region          string `json:"region"`
		State           string `json:"state"`
		City            string `json:"city"`
		District        string `json:"district"`
		Town            string `json:"town"`
		Postcode        string `json:"postcode"`
		DetailedAddress string `json:"detailed_address"`
		AdditionalInfo  string `json:"additional_info"`
		FullAddress     string `json:"full_address"`
	}
	InvoiceDetail struct {
		Name             string           `json:"name"`
		Email            string           `json:"email"`
		Address          string           `json:"address"`
		PhoneNumber      string           `json:"phone_number"`
		TaxID            string           `json:"tax_id"`
		AddressBreakdown AddressBreakdown `json:"address_breakdown"`
	}
	InvoiceInfoList struct {
		OrderSn       string        `json:"order_sn"`
		InvoiceType   string        `json:"invoice_type"`
		InvoiceDetail InvoiceDetail `json:"invoice_detail"`
		IsRequested   bool          `json:"is_requested"`
		Error         string        `json:"error"`
	}
)

type OrderServiceOp struct {
	client *ShopeeClient
}

func (o *OrderServiceOp) GetOrderDetailByOrderSN(shopID uint64, token string, params GetOrderDetailParamsRequest) (*GetOrderDetailResponse, error) {
	path := "/order/get_order_detail"

	resp := new(GetOrderDetailResponse)
	err := o.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}
