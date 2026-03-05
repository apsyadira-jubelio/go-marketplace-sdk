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

	// Business params must be in val for both signature and URL query string
	val.Set("uploadId", param.UploadId)
	val.Set("blockNo", fmt.Sprintf("%d", param.BlockNo))
	val.Set("blockCount", fmt.Sprintf("%d", param.BlockCount))

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

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
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

// ---- CompleteCreateVideo (commit) ----

type VideoPart struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"eTag"`
}

type CompleteCreateVideoRequest struct {
	UploadID   string `json:"uploadId"`
	Title      string `json:"title"`
	CoverURL   string `json:"coverUrl"`
	VideoUsage string `json:"videoUsage,omitempty"`
}

type CompleteCreateVideoResponse struct {
	Code          string `json:"code"`
	VideoID       string `json:"video_id"`
	ResultMessage string `json:"result_message"`
	Message       string `json:"message"`
	Success       bool   `json:"success"`
	ResultCode    string `json:"result_code"`
	RequestID     string `json:"request_id"`
}

// CompleteCreateVideoRaw calls /media/video/block/commit to finalize an upload.
func (m *MediaService) CompleteCreateVideoRaw(ctx context.Context, req *CompleteCreateVideoRequest, parts []VideoPart) (*CompleteCreateVideoResponse, error) {
	ts := fmt.Sprintf("%d", time.Now().Unix()*1000)
	baseURL := m.client.BaseURL.String() + "rest/media/video/block/commit"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	api := strings.TrimPrefix(u.Path, "/rest")

	// URL-encode the parts JSON
	partsJSON, err := json.Marshal(parts)
	if err != nil {
		return nil, fmt.Errorf("marshal parts: %w", err)
	}

	val := u.Query()
	val.Set("sign_method", "sha256")
	val.Set("timestamp", ts)
	val.Set("app_key", m.client.appKey)
	val.Set("access_token", m.client.accessToken)
	val.Set("uploadId", req.UploadID)
	val.Set("parts", string(partsJSON))
	val.Set("title", req.Title)
	val.Set("coverUrl", req.CoverURL)

	if req.VideoUsage != "" {
		val.Set("videoUsage", req.VideoUsage)
	}

	// Compute signature
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
	val.Set("sign", strings.ToUpper(hex.EncodeToString(sig)))

	u.RawQuery = val.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("unable to read body")
	}

	log.Println("commit raw response:", string(data))

	var lazResp CompleteCreateVideoResponse
	if err := json.Unmarshal(data, &lazResp); err != nil {
		return nil, errors.New("unable to decode response")
	}

	return &lazResp, nil
}

// ---- UploadVideo (full flow) ----

// UploadVideoResponse contains the final result of the full upload flow.
type UploadVideoResponse struct {
	UploadID string
	VideoID  string
}

// UploadVideo handles the full video upload flow in a single call:
// 1. InitCreateVideo to get upload_id
// 2. Split file into blocks and upload each
// 3. CompleteCreateVideo to commit and get video_id
//
// title is the video title (required), coverUrl is the cover image URL (required).
func (m *MediaService) UploadVideo(ctx context.Context, filename, title, coverUrl string, fileData []byte) (*UploadVideoResponse, error) {
	// Step 1: Init
	initResp, err := m.InitCreateVideo(ctx, &InitCreateVideoParameter{
		FileName:  filename,
		FileBytes: int64(len(fileData)),
	})
	if err != nil {
		return nil, fmt.Errorf("init create video: %w", err)
	}

	if initResp.UploadID == "" {
		return nil, fmt.Errorf("init create video returned empty upload_id")
	}

	// Step 2: Split and upload blocks, collect eTags
	blocks := SplitFileToBlocks(fileData, MaxBlockSizeBytes)
	parts := make([]VideoPart, 0, len(blocks))

	for i, block := range blocks {
		blockResp, err := m.UploadVideoBlockRaw(ctx, filename, &UploadVideoBlockRequest{
			UploadId:   initResp.UploadID,
			BlockNo:    i,
			BlockCount: len(blocks),
			File:       block,
		})
		if err != nil {
			return nil, fmt.Errorf("upload block %d/%d: %w", i+1, len(blocks), err)
		}

		parts = append(parts, VideoPart{
			PartNumber: i + 1,
			ETag:       blockResp.ETag,
		})
	}

	// Step 3: Commit
	commitResp, err := m.CompleteCreateVideoRaw(ctx, &CompleteCreateVideoRequest{
		UploadID: initResp.UploadID,
		Title:    title,
		CoverURL: coverUrl,
	}, parts)

	if err != nil {
		return nil, fmt.Errorf("complete create video: %w", err)
	}

	return &UploadVideoResponse{
		UploadID: initResp.UploadID,
		VideoID:  commitResp.VideoID,
	}, nil
}
