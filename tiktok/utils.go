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
