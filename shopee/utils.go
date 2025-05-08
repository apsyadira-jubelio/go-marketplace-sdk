package shopee

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type UtilService interface {
	Sign(string) (string, int64, error)
}

type UtilServiceOp struct {
	client *ShopeeClient
}

func (s *UtilServiceOp) Sign(plainText string) (string, int64, error) {
	ts := time.Now().Unix()
	baseStr := fmt.Sprintf("%d%s%d", s.client.appConfig.PartnerID, plainText, ts)
	h := hmac.New(sha256.New, []byte(s.client.appConfig.PartnerKey))
	h.Write([]byte(baseStr))
	result := hex.EncodeToString(h.Sum(nil))
	return result, ts, nil
}

func StructToMap(in interface{}) (map[string]interface{}, error) {
	byts, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("error to perpare request body: %s", err)
	}
	var res map[string]interface{}
	if err := json.Unmarshal(byts, &res); err != nil {
		return nil, fmt.Errorf("error to perpare request body 1: %s", err)
	}

	if v, ok := in.(ReadMessageRequest); ok {
		dec := json.NewDecoder(bytes.NewReader([]byte(v.ConversationID)))
		dec.UseNumber()
		if err := dec.Decode(&v.ConversationID); err != nil {
			log.Println("error decode")
			return nil, err
		}

		res["conversation_id"], _ = v.ConversationID.Int64()
	}

	if v, ok := in.(UnreadMessageRequest); ok {
		dec := json.NewDecoder(bytes.NewReader([]byte(v.ConversationID)))
		dec.UseNumber()
		if err := dec.Decode(&v.ConversationID); err != nil {
			log.Println("error decode")
			return nil, err
		}

		res["conversation_id"], _ = v.ConversationID.Int64()
	}

	if v, ok := in.(SendMessageRequest); ok {
		dec := json.NewDecoder(bytes.NewReader([]byte(v.ToID)))
		dec.UseNumber()
		if err := dec.Decode(&v.ToID); err != nil {
			log.Println("error decode")
			return nil, err
		}

		// required
		res["to_id"], _ = v.ToID.Int64()

		if vItem, err := v.Content.ItemID.Float64(); err == nil && vItem > 0 {
			decItemID := json.NewDecoder(bytes.NewReader([]byte(v.Content.ItemID)))
			decItemID.UseNumber()
			if err := decItemID.Decode(&v.Content.ItemID); err != nil {
				log.Println("error decode")
				return nil, err
			}

			res["item_id"], _ = v.Content.ItemID.Int64()
		}
	}

	return res, nil
}
