package tiktok

import "fmt"

type FulfillmentService interface {
	GetTracking(orderID string) (*GetTrackingResponse, error)
}

type FulfillmentServiceOp struct {
	client *TiktokClient
}

type GetTrackingParamRequest struct {
	OrderID string `url:"order_id"`
}

type GetTrackingResponse struct {
	BaseResponse
	Data *TrackingData `json:"data"`
}

type Tracking struct {
	Description      string `json:"description"`
	UpdateTimeMillis int64  `json:"update_time_millis"`
}
type TrackingData struct {
	Tracking []Tracking `json:"tracking"`
}

// /fulfillment/202309/orders/{order_id}/tracking
func (p *FulfillmentServiceOp) GetTracking(orderID string) (*GetTrackingResponse, error) {
	path := fmt.Sprintf("/fulfillment/%s/orders/%s/tracking", p.client.appConfig.Version, orderID)

	resp := new(GetTrackingResponse)
	err := p.client.Get(path, resp, nil)

	return resp, err
}
