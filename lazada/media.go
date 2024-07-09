package lazada

import (
	"context"
	"encoding/json"
)

// The Media Service deals with any methods under the "Media Center API" category of the open platform
type MediaService service

type GetVideoParameter struct {
	VideoID string `url:"videoId" json:"videoId"`
}

type GetVideoResponse struct {
	CoverURL      string `json:"cover_url"`
	VideoURL      string `json:"video_url"`
	Code          string `json:"code"`
	ResultMessage string `json:"result_message"`
	Success       string `json:"success"`
	ResultCode    string `json:"result_code"`
	State         string `json:"state"`
	Title         string `json:"title"`
	RequestID     string `json:"request_id"`
}

func (o *MediaService) GetVideo(ctx context.Context, opts *GetVideoParameter) (res *GetVideoResponse, err error) {
	u, err := addOptions(ApiNames["GetVideo"], opts)
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

	json.Unmarshal(jsonData, &res)

	return res, nil
}
