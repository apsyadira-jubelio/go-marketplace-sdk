package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/lazada"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	appKey := os.Getenv("LAZADA_APP_KEY")
	appSecret := os.Getenv("LAZADA_APP_SECRET")
	client := lazada.NewClient(appKey, appSecret, lazada.Indonesia)
	client.NewTokenClient(os.Getenv("LAZADA_TOKEN"))

	videoPath := "./test.mp4"

	// Step 1: Extract thumbnail from video using ffmpeg
	thumbnailPath, err := extractThumbnail(videoPath)
	if err != nil {
		log.Fatal("Failed to extract thumbnail:", err)
	}
	defer os.Remove(thumbnailPath)
	log.Printf("Thumbnail extracted: %s\n", thumbnailPath)

	// Step 2: Upload thumbnail to file service to get a public URL
	coverUrl, err := uploadToFileService(thumbnailPath, os.Getenv("STORAGE_TOKEN"))
	if err != nil {
		log.Fatal("Failed to upload thumbnail:", err)
	}
	log.Printf("Cover URL: %s\n", coverUrl)

	// Step 3: Read video file
	fileData, err := os.ReadFile(videoPath)
	if err != nil {
		log.Fatal(err)
	}

	// Step 4: Upload video (init + blocks + commit) in one call
	resp, err := client.Media.UploadVideo(context.Background(), filepath.Base(videoPath), "Test Video", coverUrl, fileData)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Upload success! Upload ID: %s, Video ID: %s\n", resp.UploadID, resp.VideoID)
}

// extractThumbnail uses ffmpeg to extract a frame from a video as a JPEG.
func extractThumbnail(videoPath string) (string, error) {
	outputPath := fmt.Sprintf("%s_thumb.jpg", videoPath)

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-q:v", "2",
		"-y",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %w\noutput: %s", err, string(output))
	}

	return outputPath, nil
}

// uploadToFileService uploads an image to the storage service and returns the public URL.
func uploadToFileService(filePath, bearerToken string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Detect MIME type from extension
	ext := filepath.Ext(filePath)
	mimeType := "image/jpeg"
	if ext == ".png" {
		mimeType = "image/png"
	}

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, filepath.Base(filePath)))
	partHeader.Set("Content-Type", mimeType)

	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("write file data: %w", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "https://chat-api.qm-staging-k8s.jubelio.io/storage/upload", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse response to get the URL
	var result struct {
		URL  string      `json:"url"`
		Data interface{} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w\nraw: %s", err, string(respBody))
	}

	// Try top-level url first
	imageURL := result.URL
	if imageURL == "" {
		// Try data field - could be string or nested struct
		switch v := result.Data.(type) {
		case string:
			imageURL = v
		case map[string]interface{}:
			if u, ok := v["url"].(string); ok {
				imageURL = u
			}
		}
	}
	if imageURL == "" {
		return "", fmt.Errorf("no url in response: %s", string(respBody))
	}

	return imageURL, nil
}
