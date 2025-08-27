package lazada

import (
	"context"
	"encoding/json"
)

// LogisticService handles operations related to logistics.
type LogisticService service

// GetOrderTraceParams represents the parameters for the GetOrderTrace API call.
type GetOrderTraceParams struct {
	OrderID   string `url:"order_id"`
	PackageID string `url:"package_id,omitempty"`
	Locale    string `url:"locale,omitempty"`
}

// GetOrderTraceResponse represents the response from the GetOrderTrace API call.
type GetOrderTraceResponse struct {
	BaseResponse
	Data LogisticTraceResponseData `json:"data"`
}

// LogisticTraceResponseData contains the details of the logistic trace.
type LogisticTraceResponseData struct {
	MailNo         string          `json:"mail_no"`
	OfcName        string          `json:"ofc_name"`
	OfcPhone       string          `json:"ofc_phone"`
	LogisticEvents []LogisticEvent `json:"logistic_event_list"`
}

// LogisticEvent represents a single event in the logistic trace.
type LogisticEvent struct {
	EventTime   string `json:"event_time"`
	EventDesc   string `json:"event_desc"`
	EventCode   string `json:"event_code"`
	OfcName     string `json:"ofc_name"`
	OfcPostCode string `json:"ofc_post_code"`
	OfcCity     string `json:"ofc_city"`
	OfcProvince string `json:"ofc_province"`
	OfcCountry  string `json:"ofc_country"`
	Signatory   string `json:"signatory"`
	EventWeight int    `json:"event_weight"`
}

// GetOrderTrace retrieves the logistic trace information for a given order.
func (s *LogisticService) GetOrderTrace(ctx context.Context, params *GetOrderTraceParams) (*GetOrderTraceResponse, error) {
	u, err := addOptions(ApiNames["OrderTrace"], params)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	resp := &GetOrderTraceResponse{}
	if err := json.Unmarshal(jsonData, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
