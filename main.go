package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/lazada"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	appKey := "117532"
	appSecret := "G5VwB0wyhk3XQEsklCfmHSF2kP2luEqS"
	// playground
	client := lazada.NewClient(appKey, appSecret, lazada.Indonesia)
	client.NewTokenClient("50000501928fIpdsqf7kcs6N0UHgEJjYGvTBfbCv15c4549cPcxtyqsRXpaTjqIj")

	initFileBytes, err := os.ReadFile("./test.mp4")
	if err != nil {
		log.Fatal(err)
	}

	// resp, err := client.Media.InitCreateVideo(context.Background(), &lazada.InitCreateVideoParameter{
	// 	FileName:  "test.mp4",
	// 	FileBytes: 3145728,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := client.Media.UploadVideoBlockRaw(context.Background(), "test.mp4", &lazada.UploadVideoBlockRequest{
		UploadId:   "DD555F3DDB5844F8A7E8E0F5A7B2647A",
		BlockNo:    0,
		BlockCount: 1,
		File:       initFileBytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	// cara kedua, chunk video jadi beberapa bagian sesuai sama: https://open.lazada.com/apps/doc/api?path=%2Fmedia%2Fvideo%2Fblock%2Fupload
	// blocks := lazada.SplitFileToBlocks(initFileBytes, lazada.MaxBlockSizeBytes)
	// for i, block := range blocks {
	// 	resp, err := client.Media.UploadVideoBlockRaw(context.Background(), "test.mp4", &lazada.UploadVideoBlockRequest{
	// 		UploadId:   "DD555F3DDB5844F8A7E8E0F5A7B2647A",
	// 		BlockNo:    i,
	// 		BlockCount: len(blocks),
	// 		File:       block,
	// 	})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	writeJSONFile(resp, fmt.Sprintf("response-upload-block-%d", i))
	// }

	spew.Dump(resp)
}

func writeJSONFile(response any, filename string) error {
	// Create a new JSON file
	file, err := os.Create(fmt.Sprintf("%s.json", filename))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Encode response data to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Indent the JSON for readability (optional)
	err = encoder.Encode(response)
	if err != nil {
		fmt.Println("Error encoding JSON to file:", err)
		return err
	}

	fmt.Printf("Response written to %s.json", filename)
	return nil
}
