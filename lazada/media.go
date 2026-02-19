package lazada

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
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

func (m *MediaService) GetVideo(ctx context.Context, opts *GetVideoParameter) (res *GetVideoResponse, err error) {
	u, err := addOptions(ApiNames["GetVideo"], opts)
	if err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(ctx, req, nil)
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

type InitCreateVideoParameter struct {
	FileName  string `url:"fileName" json:"fileName"`
	FileBytes int64  `url:"fileBytes" json:"fileBytes"` // size of file
}

type InitCreateVideoResponse struct {
	UploadID      string `json:"upload_id"`
	Code          string `json:"code"`
	ResultMessage string `json:"result_message"`
	Success       bool   `json:"success"`
	ResultCode    string `json:"result_code"`
	Title         string `json:"title"`
	RequestID     string `json:"request_id"`
}

func (m *MediaService) InitCreateVideo(ctx context.Context, opts *InitCreateVideoParameter) (res *InitCreateVideoResponse, err error) {
	u, err := addOptions(ApiNames["InitCreateVideo"], opts)
	if err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(ctx, req, nil)
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

type UploadVideoBlockRequest struct {
	UploadId   string `url:"uploadId"`
	BlockNo    int    `url:"blockNo"`
	BlockCount int    `url:"blockCount"`
	File       []byte `url:"file"`
}

type UploadVideoBlockResponse struct {
	UploadID      string `json:"upload_id"`
	Code          string `json:"code"`
	ResultMessage string `json:"result_message"`
	Message       string `json:"message"`
	Success       bool   `json:"success"`
	ResultCode    string `json:"result_code"`
	ETag          string `json:"e_tag"`
	RequestID     string `json:"request_id"`
}

// The API is used to upload one block of origin video file.
// The video file can split into multiple files. For example, a 8MB video file can be split into three blocks. 3MB, 3MB and 2MB.
// These three blocks can be uploaded by calling UploadVideoBlock three times.
func (m *MediaService) UploadVideoBlock(ctx context.Context, opts *UploadVideoBlockRequest) (res *UploadVideoBlockResponse, err error) {
	u, err := addOptions(ApiNames["UploadVideoBlock"], opts)
	if err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}
	log.Println("check url:", req.URL)

	var buf bytes.Buffer
	_, err = m.client.Do(ctx, req, &buf)
	if err != nil {
		return nil, err
	}

	t := &UploadVideoBlockResponse{}
	if err := json.NewDecoder(&buf).Decode(t); err != nil {
		return nil, errors.New("cant unmarshal rep upload block")
	}

	return res, nil
}

func (m *MediaService) UploadVideoBlockRaw(ctx context.Context, filename string, param *UploadVideoBlockRequest) (*UploadVideoBlockResponse, error) {
	ts := fmt.Sprintf("%d", time.Now().Unix()*1000)
	baseURL := "https://api.lazada.co.id/rest/media/video/block/upload"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	api := strings.TrimPrefix(u.Path, "/rest")
	val := u.Query()

	val.Set("sign_method", "sha256")
	val.Set("timestamp", ts)
	val.Set("app_key", m.client.appKey)
	val.Set("access_token", m.client.accessToken)
	var buf bytes.Buffer
	buf.WriteString(api)
	keys := make([]string, 0, len(val))
	for k := range val {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vs := val[k]
		keyEscaped := url.QueryEscape(k)

		for _, v := range vs {
			buf.WriteString(keyEscaped)
			buf.WriteString(v)
		}
	}

	signer := hmac.New(sha256.New, []byte(m.client.secret))
	signer.Write(buf.Bytes())
	sig := signer.Sum(nil)

	sign := strings.ToUpper(hex.EncodeToString(sig))
	val.Set("sign", sign)

	finalURL := fmt.Sprintf("%s?uploadId=%s&blockNo=%d&blockCount=%d&sign_method=sha256&sign=%s&timestamp=%s&app_key=%s&access_token=%s", u.String(), param.UploadId, param.BlockNo, param.BlockCount, sign, ts, m.client.appKey, m.client.accessToken)
	log.Println("finalURL: ", finalURL)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(param.File); err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", finalURL, body)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = val.Encode()

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("unable to read body")
	}

	var lazResp UploadVideoBlockResponse
	decErr := json.Unmarshal(data, &lazResp)
	if decErr != nil {
		return nil, errors.New("unable to decode response")
	}

	return &lazResp, nil

}
