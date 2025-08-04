package lazada

import (
	"context"
	"encoding/json"
)

type VoucherService service

type GetVouchersParam struct {
	CurPage     string `url:"cur_page"`
	VoucherType string `url:"voucher_type"`
	PageSize    string `url:"page_size,omitempty"`
	Name        string `url:"name,omitempty"`
	Status      string `url:"status,omitempty"`
}

type GetVoucherResponse struct {
	BaseResponse
	Data any `json:"data"`
}

func (o *VoucherService) GetVouchers(ctx context.Context, opts *GetVouchersParam) (res *GetVoucherResponse, err error) {
	u, err := addOptions(ApiNames["GetVouchers"], opts)
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
