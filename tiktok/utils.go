package tiktok

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type UtilService interface {
	Sign(string) (string, int64, error)
}

type UtilServiceOp struct {
	client *TiktokClient
}

func (s *UtilServiceOp) Sign(plainText string) (string, int64, error) {
	ts := time.Now().Unix()
	baseStr := fmt.Sprintf("%s%s%d", s.client.appConfig.AppKey, plainText, ts)
	h := hmac.New(sha256.New, []byte(s.client.appConfig.AppSecret))
	h.Write([]byte(baseStr))
	result := hex.EncodeToString(h.Sum(nil))
	return result, ts, nil
}

const chunkSize = 10 * 1024 * 1024
const minTailSize = 5 * 1024 * 1024

func CalcChunkCount(size int64) int {
	if size <= 0 {
		return 1
	}
	n := int(size / chunkSize)
	rem := size % chunkSize
	if rem > 0 {
		n++
	}
	if rem > 0 && rem < minTailSize && n > 1 {
		n--
	}
	if n <= 0 {
		n = 1
	}
	return n
}
