package shopee

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ChatService interface {
	GetMessage(shopID uint64, token string, params GetMessageParamsRequest) (*GetMessageResponse, error)
	GetConversationList(shopID uint64, token string, params GetConversationParamsRequest) (*GetConversationResponse, error)
	GetOneConversation(shopID uint64, token string, params GetMessageParamsRequest) (*GetDetailConversation, error)
	SendMessage(shopID uint64, token string, request SendMessageRequest) (*GetSendMessageResponse, error)
	UploadImage(shopID uint64, token string, filename string) (*UploadImageResponse, error)
	GetStickerPack() (*StickerPacksResponse, error)
	GetListStickerByPID(stickerPackageID string) (*ListStickerByPID, error)
	GetStickerByPIDAndSID(stickerPackageID, stickerID string) string
	ReadConversation(shopID uint64, token string, params ReadMessageRequest) (*ReadMessageResponse, error)
	UnreadConversation(shopID uint64, token string, request UnreadMessageRequest) (*UnreadMessageResponse, error)
}

type GetMessageParamsRequest struct {
	Offset         string  `url:"offset,omitempty"`
	PageSize       int     `url:"page_size,omitempty"`
	ConversationID int64   `url:"conversation_id,omitempty"`
	MessageIdList  []int64 `url:"message_id_list,ommitempty"`
}

type GetMessageResponse struct {
	BaseResponse

	Response GetMessageDataResponse `json:"response"`
}

type GetMessageDataResponse struct {
	MessagesList []Messages `json:"messages"`
}

type Messages struct {
	MessageID        string         `json:"message_id"`
	MessageType      string         `json:"message_type"`
	FromID           int64          `json:"from_id"`
	FromShopID       int64          `json:"from_shop_id"`
	ToID             int64          `json:"to_id"`
	ToShopID         int64          `json:"to_shop_id"`
	ConversationID   string         `json:"conversation_id"`
	CreatedTimeStamp int64          `json:"created_timestamp"`
	Region           string         `json:"region"`
	Status           string         `json:"status"`
	Source           string         `json:"source"`
	Content          ContentMessage `json:"content"`
	SourceContent    SourceContent  `json:"source_content,omitempty"`
	MessageOption    int            `json:"message_option"`
}

type ContentMessage struct {
	Text             string        `json:"text,omitempty"`
	Url              string        `json:"url,omitempty"`
	ThumbHeight      int           `json:"thumb_height,omitempty"`
	ThumbWidth       int           `json:"thumb_width,omitempty"`
	ThumbURL         string        `json:"thumb_url,omitempty"`
	FileServerID     int64         `json:"file_server_id,omitempty"`
	ShopID           int64         `json:"shop_id,omitempty"`
	OfferID          int           `json:"offer_id,omitempty"`
	ProductID        int           `json:"product_id,omitempty"`
	TaxValue         string        `json:"tax_value,omitempty"`
	PriceBeforeTax   string        `json:"price_before_tax,omitempty"`
	TaxApplicable    bool          `json:"tax_applicable,omitempty"`
	StickerID        string        `json:"sticker_id,omitempty"`
	StickerPackageID string        `json:"sticker_package_id,omitempty"`
	ItemID           int64         `json:"item_id,omitempty"`
	OrderID          int64         `json:"order_id,omitempty"`
	VideoURL         string        `json:"video_url,omitempty"`
	ImageURL         string        `json:"image_url,omitempty"`
	VoucherID        string        `json:"voucher_id,omitempty"`
	VoucherCode      string        `json:"voucher_code,omitempty"`
	SourceContent    SourceContent `json:"source_content,omitempty"`
}

type SourceContent struct {
	OrderSN string `json:"order_sn,omitempty"`
	ItemID  int64  `json:"item_id,omitempty"`
}

type ChatServiceOp struct {
	client *ShopeeClient
}

func (s *ChatServiceOp) GetMessage(shopID uint64, token string, params GetMessageParamsRequest) (*GetMessageResponse, error) {
	path := "/sellerchat/get_message"

	resp := new(GetMessageResponse)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}

type GetConversationParamsRequest struct {
	Direction    string `url:"direction"` // latest/older
	Type         string `url:"type"`
	NextTimeNano int64  `url:"next_timestamp_nano,omitempty"`
	PageSize     int    `url:"page_size"`
}

type GetConversationResponse struct {
	BaseResponse

	Response GetConversationDataResponse `json:"response"`
}

type GetConversationDataResponse struct {
	PageResult        ConversationPageResult `json:"page_result"`
	ConversationsList []Conversation         `json:"conversations"`
}

type ConversationPageResult struct {
	PageSize   int `json:"page_size"`
	NextCursor struct {
		NextMessageTimeNano string `json:"next_message_time_nano"`
		ConversationID      string `json:"conversation_id"`
	} `json:"next_cursor"`
	More bool `json:"more"`
}

func (s *ChatServiceOp) GetConversationList(shopID uint64, token string, params GetConversationParamsRequest) (*GetConversationResponse, error) {
	path := "/sellerchat/get_conversation_list"

	opt := GetConversationParamsRequest{
		PageSize:     params.PageSize,
		Direction:    params.Direction,
		Type:         params.Type,
		NextTimeNano: params.NextTimeNano,
	}

	resp := new(GetConversationResponse)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, opt)
	return resp, err
}

type GetSendMessageResponse struct {
	BaseResponse

	Response GetSendMessageDataResponse `json:"response"`
}

type GetSendMessageDataResponse struct {
	MessageID   string `json:"message_id"`
	ToID        int    `json:"to_id"`
	MessageType string `json:"message_type"`
	Content     struct {
		Text string `json:"text"`
	} `json:"content"`
	ConversationID   int64 `json:"conversation_id"`
	CreatedTimestamp int   `json:"created_timestamp"`
	MessageOption    int   `json:"message_option"`
}

type SendMessageRequest struct {
	ToID           json.Number        `json:"to_id"`
	MessageType    string             `json:"message_type"`
	Content        ContentSendMessage `json:"content"`
	BusinessType   int8               `json:"business_type,omitempty"`
	ConversationID int64              `json:"conversation_id,omitempty"`
}

type ContentSendMessage struct {
	Text string `json:"text,omitempty"`
	// sticker
	StickerID        string `json:"sticker_id,omitempty"`
	StickerPackageID string `json:"sticker_package_id,omitempty"`
	// image
	ImageURL string `json:"image_url,omitempty"`
	// voucher
	VoucherID   string `json:"voucher_id,omitempty"`
	VoucherCode string `json:"voucher_code,omitempty"`
	// video
	Vid             int32  `json:"vid,omitempty"`
	VideoURL        string `json:"video_url,omitempty"`
	DurationSeconds int32  `json:"duration_seconds,omitempty"`
	// product
	ItemID json.Number `json:"item_id,omitempty"`
	// order
	OrderSN string `json:"order_sn,omitempty"`
}

func (s *ChatServiceOp) SendMessage(shopID uint64, token string, request SendMessageRequest) (*GetSendMessageResponse, error) {
	path := "/sellerchat/send_message"
	resp := new(GetSendMessageResponse)
	req, err := StructToMap(request)
	if err != nil {
		return nil, err
	}

	err = s.client.WithShop(uint64(shopID), token).Post(path, req, resp)
	return resp, err
}

type UploadImageResponse struct {
	BaseResponse

	Response UploadImageDataResponse `json:"response"`
}

type UploadImageDataResponse struct {
	FileServerID int    `json:"file_server_id"`
	ThumbHeight  int    `json:"thumb_height"`
	ThumbWidth   int    `json:"thumb_width"`
	Thumbnail    string `json:"thumbnail"`
	URL          string `json:"url"`
	URLHash      string `json:"url_hash"`
}

func (s *ChatServiceOp) UploadImage(shopID uint64, token string, filename string) (*UploadImageResponse, error) {
	path := "/sellerchat/upload_image"

	resp := new(UploadImageResponse)
	err := s.client.WithShop(uint64(shopID), token).Upload(path, "file", filename, resp)
	return resp, err
}

type GetDetailConversation struct {
	BaseResponse
	Response Conversation `json:"response"`
}

type Conversation struct {
	ConversationID       string `json:"conversation_id"`
	ToID                 int    `json:"to_id"`
	ToName               string `json:"to_name"`
	ToAvatar             string `json:"to_avatar"`
	ShopID               int    `json:"shop_id"`
	UnreadCount          int    `json:"unread_count"`
	Pinned               bool   `json:"pinned"`
	Mute                 bool   `json:"mute"`
	LastReadMessageID    string `json:"last_read_message_id"`
	LatestMessageID      string `json:"latest_message_id"`
	LatestMessageType    string `json:"latest_message_type"`
	LatestMessageContent struct {
		Text string `json:"text"`
	} `json:"latest_message_content"`
	LatestMessageFromID      int    `json:"latest_message_from_id"`
	LastMessageTimestamp     int64  `json:"last_message_timestamp"`
	LastMessageOption        int    `json:"last_message_option"`
	MaxGeneralOptionHideTime string `json:"max_general_option_hide_time"`
}

// Use GetMessageParamsRequest, need param convesation_id
func (s *ChatServiceOp) GetOneConversation(shopID uint64, token string, params GetMessageParamsRequest) (*GetDetailConversation, error) {
	path := "/sellerchat/get_one_conversation"

	resp := new(GetDetailConversation)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}

type StickerPacksResponse struct {
	Packs []Packs `json:"packs"`
}
type Packs struct {
	Md5 string   `json:"md5"`
	Pid string   `json:"pid"`
	Reg []string `json:"reg"`
}

func (s *ChatServiceOp) GetStickerPack() (*StickerPacksResponse, error) {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 30 * time.Second,
	}

	url := "https://deo.shopeemobile.com/shopee/shopee-sticker-live-id/manifest.json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading body, cause: %+v\n", err)
			return nil, err
		}
		log.Println("response:", map[string]interface{}{
			"body": string(respBytes),
			"code": resp.StatusCode,
		})
		return nil, errors.New("response not ok")
	}
	defer resp.Body.Close()

	var response StickerPacksResponse
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body, cause: %+v\n", err)
		return nil, err
	}
	if err := json.Unmarshal(respBytes, &response); err != nil {
		log.Printf("Error unmarshaling response, cause: %+v\n", err)
		log.Println("Response:", string(respBytes))
		return nil, err
	}

	return &response, nil
}

type ListStickerByPID struct {
	AutoDownload bool       `json:"auto_download"`
	Locales      []string   `json:"locales"`
	Size         []int      `json:"size"`
	Stickers     []Stickers `json:"stickers"`
}
type Stickers struct {
	Ext  string   `json:"ext"`
	Name []string `json:"name"`
	Sid  string   `json:"sid"`
}

func (s *ChatServiceOp) GetListStickerByPID(stickerPackageID string) (*ListStickerByPID, error) {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 30 * time.Second,
	}

	url := fmt.Sprintf("https://deo.shopeemobile.com/shopee/shopee-sticker-live-id/packs/%s/%s.json", stickerPackageID, stickerPackageID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading body, cause: %+v\n", err)
			return nil, err
		}
		log.Println("response:", map[string]interface{}{
			"body": string(respBytes),
			"code": resp.StatusCode,
		})
		return nil, errors.New("response not ok")
	}
	defer resp.Body.Close()

	var response ListStickerByPID
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body, cause: %+v\n", err)
		return nil, err
	}
	if err := json.Unmarshal(respBytes, &response); err != nil {
		log.Printf("Error unmarshaling response, cause: %+v\n", err)
		log.Println("Response:", string(respBytes))
		return nil, err
	}

	return &response, nil
}

func (s *ChatServiceOp) GetStickerByPIDAndSID(stickerPackageID, stickerID string) string {
	filename := fmt.Sprintf("https://deo.shopeemobile.com/shopee/shopee-sticker-live-id/packs/%s/%s@1x.png", stickerPackageID, stickerID)
	return filename
}

type ReadMessageResponse struct {
	BaseResponse
	Response any `json:"response"`
}

type ReadMessageRequest struct {
	ConversationID    json.Number `json:"conversation_id"`
	LastReadMessageID string      `json:"last_read_message_id"`
	BusinessType      int32       `json:"business_type,omitempty"`
}

// ReadConversation for read message from buyer by last_read_message_id from response get message
func (s *ChatServiceOp) ReadConversation(shopID uint64, token string, request ReadMessageRequest) (*ReadMessageResponse, error) {
	path := "/sellerchat/read_conversation"
	req, err := StructToMap(request)
	if err != nil {
		return nil, err
	}

	resp := new(ReadMessageResponse)
	err = s.client.WithShop(uint64(shopID), token).Post(path, req, resp)
	return resp, err
}

type UnreadMessageResponse struct {
	BaseResponse
	Response any `json:"response"`
}

type UnreadMessageRequest struct {
	ConversationID json.Number `json:"conversation_id"`
	BusinessType   int32       `json:"business_type,omitempty"` // not required, 0 is for seller buyer chat, 11 is for seller affiliate chat
}

// UnreadConversation for mark a conversation from buyer as unread
func (s *ChatServiceOp) UnreadConversation(shopID uint64, token string, request UnreadMessageRequest) (*UnreadMessageResponse, error) {
	path := "/sellerchat/unread_conversation"
	resp := new(UnreadMessageResponse)
	req, err := StructToMap(request)
	if err != nil {
		return nil, err
	}
	err = s.client.WithShop(uint64(shopID), token).Post(path, req, resp)
	return resp, err
}
