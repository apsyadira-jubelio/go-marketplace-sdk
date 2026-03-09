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
	"os/exec"
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
	Success       bool   `json:"success"`
	ResultCode    string `json:"result_code"`
	State         string `json:"state"`
	Title         string `json:"title"`
	ErrCode       string `json:"err_code,omitempty"`
	ErrMessage    string `json:"err_message,omitempty"`
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
// 1. Split file into blocks and upload each
// 2. CompleteCreateVideo to commit and get video_id
//
// title is the video title (required), coverUrl is the cover image URL (required).
func (m *MediaService) UploadVideo(ctx context.Context, filename, title, coverUrl, uploadID string, fileData []byte) (*UploadVideoResponse, error) {
	// Step 1: Validate uploadID
	if uploadID == "" {
		return nil, fmt.Errorf("init create video returned empty upload_id")
	}

	// Step 2: Split and upload blocks, collect eTags
	blocks := SplitFileToBlocks(fileData, MaxBlockSizeBytes)
	parts := make([]VideoPart, 0, len(blocks))

	for i, block := range blocks {
		blockResp, err := m.UploadVideoBlockRaw(ctx, filename, &UploadVideoBlockRequest{
			UploadId:   uploadID,
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
		UploadID: uploadID,
		Title:    title,
		CoverURL: coverUrl,
	}, parts)

	if err != nil {
		return nil, fmt.Errorf("complete create video: %w", err)
	}

	return &UploadVideoResponse{
		UploadID: uploadID,
		VideoID:  commitResp.VideoID,
	}, nil
}

type ThumbnailOptions struct {
	TimeOffset string
	Quality    int
}

// ExtractVideoThumbnailToBytes extracts a frame from video bytes without writing to disk.
func (m *MediaService) ExtractVideoThumbnailToBytes(videoData []byte, opts *ThumbnailOptions) ([]byte, error) {
	timeOffset := "00:00:01"
	quality := 2
	if opts != nil {
		if opts.TimeOffset != "" {
			timeOffset = opts.TimeOffset
		}
		if opts.Quality > 0 {
			quality = opts.Quality
		}
	}

	// Input from pipe (stdin) with pipe:0, output to pipe (stdout) with image2pipe format
	cmd := exec.Command("ffmpeg",
		"-ss", timeOffset,
		"-i", "pipe:0", // read from stdin
		"-frames:v", "1",
		"-q:v", fmt.Sprintf("%d", quality),
		"-vf", "scale=800:-1",
		"-f", "image2pipe", // output as raw image stream
		"-vcodec", "mjpeg", // encode as JPEG
		"pipe:1", // send to stdout
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(videoData)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w\n%s", err, stderr.String())
	}
	if stdout.Len() == 0 {
		return nil, errors.New("ffmpeg: output is empty")
	}

	return stdout.Bytes(), nil
}
