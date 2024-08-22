package lazada

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// The Chat Service deals with any methods under the "Instant Messaging" category of the open platform
type ChatService service

// Base Response to return Session List
type GetSessionListResponse struct {
	BaseResponse
	SessionListResponseData `json:"data"`
}

type SessionListData struct {
	Summary         string   `json:"summary"`
	UnreadCount     int      `json:"unread_count"`
	LastMessageID   string   `json:"last_message_id"`
	HeadURL         string   `json:"head_url"`
	SelfPosition    int      `json:"self_position"`
	SiteID          string   `json:"site_id"`
	LastMessageTime int64    `json:"last_message_time"`
	SessionID       string   `json:"session_id"`
	BuyerID         int64    `json:"buyer_id"`
	Title           string   `json:"title"`
	ToPosition      int      `json:"to_position"`
	Tags            []string `json:"tags"`
}

type SessionListResponseData struct {
	SessionList   []SessionListData `json:"session_list"`
	NextStartTime int64             `json:"next_start_time"`
	HasMore       bool              `json:"has_more"`
	LastSessionID string            `json:"last_session_id"`
}

type SessionListQuery struct {
	LastSessionID string `url:"last_session_id,omitempty"`
	StartTime     int64  `url:"start_time"`
	PageSize      int    `url:"page_size"`
}

func epochTimeOneMonthAgo() int64 {
	// Get the current time.
	now := time.Now()
	// Subtract one month.
	oneMonthAgo := now.AddDate(0, -1, 0)
	// Get the Unix timestamp in milliseconds.
	timestamp := oneMonthAgo.UnixNano() / int64(time.Millisecond)
	return timestamp
}

// GetSessionList retrieves a list of sessions based on the provided query options.
// It returns a GetSessionListResponse containing the session list and any errors encountered.
// If the opts parameter is nil, default options will be used with a page size of 20 and start time set to one month ago.
func (m *ChatService) GetSessionList(ctx context.Context, opts *SessionListQuery) (*GetSessionListResponse, error) {
	if opts == nil {
		opts = &SessionListQuery{
			PageSize:  20,
			StartTime: epochTimeOneMonthAgo(),
		}
	}

	u, err := addOptions(ApiNames["GetSessionList"], opts)
	if err != nil {
		return nil, err
	}

	req, err := m.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	res, err := m.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	resp := &GetSessionListResponse{}
	json.Unmarshal([]byte(jsonData), &resp)

	return resp, nil
}

// A session detail object returned from the open platform
type GetSessionDetailResponse struct {
	BaseResponse
	SessionListData `json:"data"`
}

// GetSessionDetail returns a list of session detail in the region set
// sessionID is required
func (m *ChatService) GetSessionDetail(ctx context.Context, sessionID string) (*GetSessionDetailResponse, error) {
	req, err := m.client.NewRequest("GET", fmt.Sprintf("%s?session_id=%s", ApiNames["GetSessionDetail"], sessionID), nil)
	if err != nil {
		return nil, err
	}

	res, err := m.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	resp := &GetSessionDetailResponse{}
	json.Unmarshal([]byte(jsonData), &resp)

	return resp, nil
}

// A messages object returned from the open platform
type GetMessageResponse struct {
	BaseResponse

	Data GetMessageResponseData `json:"data"`
}

type GetMessageResponseData struct {
	LastMessageID string             `json:"last_message_id"`
	MessageList   []MessagesListData `json:"message_list"`
	NextStartTime int                `json:"next_start_time"`
	HasMore       bool               `json:"has_more"`
}

type MessagesListData struct {
	FromAccountType int    `json:"from_account_type"`
	ProcessMsg      string `json:"process_msg"`
	SessionID       string `json:"session_id"`
	MessageID       string `json:"message_id"`
	Type            string `json:"type"`
	Content         string `json:"content"`
	ToAccountID     string `json:"to_account_id"`
	SendTime        int    `json:"send_time"`
	AutoReply       string `json:"auto_reply"`
	ToAccountType   int    `json:"to_account_type"`
	SiteID          string `json:"site_id"`
	TemplateID      int    `json:"template_id"`
	FromAccountID   string `json:"from_account_id"`
	Status          string `json:"status"`
}

type MessageQueryParams struct {
	SessionID     string `url:"session_id"`
	StartTime     int64  `url:"start_time"`
	PageSize      int    `url:"page_size"`
	LastMessageID string `url:"last_message_id,omitempty"`
}

// GetMessageList lets you retrieve all message list in a specific session
func (m *ChatService) GetMessageList(ctx context.Context, opts *MessageQueryParams) (res *GetMessageResponse, err error) {
	u, err := addOptions(ApiNames["GetMessageList"], opts)
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

	json.Unmarshal([]byte(jsonData), &res)

	return res, nil
}

// MessageRecallParams is a struct that holds parameters for the MessageRecall function.
type MessageRecallParams struct {
	SessionID string `url:"session_id"`
	MessageID string `url:"message_id"`
}

// MessageRecall is a method on the ChatService struct. It sends a request to the server to recall a specific message.
// opts: A pointer to a MessageRecallParams struct containing the parameters for the MessageRecall function.
// The function returns a pointer to a BaseResponse struct containing the server's response, and an error, if there is one.
func (m *ChatService) MessageRecall(ctx context.Context, opts *MessageRecallParams) (res *BaseResponse, err error) {
	u, err := addOptions(ApiNames["MessageRecall"], opts)
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

	json.Unmarshal([]byte(jsonData), &res)

	return
}

// OpenSessionResposne is a struct that holds the base response along with the unique identifier of a chat session.
type OpenSessionResposne struct {
	BaseResponse

	SessionID string `json:"session_id"`
}

// OpenSession is a method on the ChatService struct. It sends a request to the server to open a chat session.
// orderID: A string representing the unique identifier of an order.
// The function returns a pointer to an OpenSessionResposne struct containing the server's response, and an error, if there is one.
func (m *ChatService) OpenSession(ctx context.Context, orderID string) (res *OpenSessionResposne, err error) {
	req, err := m.client.NewRequest("GET", fmt.Sprintf("%s?order_id=%s", ApiNames["OpenSession"], orderID), nil)
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

	json.Unmarshal([]byte(jsonData), &res)

	return
}

type ReadSessionResponse BaseResponse

// ReadSessionParams is a struct that holds parameters for reading a session in a chat service.
type ReadSessionParams struct {
	SessionID         string `url:"session_id"`
	LastReadMessageID string `url:"last_read_message_id"`
}

// ReadSession is a method on the ChatService struct. It sends a request to the server to read message of chat session.
// opts: A pointer to a ReadSessionParams struct containing the parameters for the ReadSession function.
// The function returns a pointer to an OpenSessionResposne struct containing the server's response, and an error, if there is one.
func (m *ChatService) ReadSession(ctx context.Context, opts ReadSessionParams) (res *ReadSessionResponse, err error) {
	u, err := addOptions(ApiNames["ReadSession"], opts)
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

	json.Unmarshal([]byte(jsonData), &res)

	return res, nil
}

type SendMessageResponse struct {
	BaseResponse
	Data struct {
		MessageID   string `json:"message_id"`
		TemplateID  int    `json:"template_id"`
		CurrentTime int64  `json:"current_time"`
	} `json:"data"`
}

type SendMessageParams struct {
	SessionID   string `url:"session_id" json:"session_id"`
	TemplateID  int    `url:"template_id" json:"template_id"`
	Txt         string `url:"txt,omitempty" json:"txt"`
	ImgUrl      string `url:"img_url,omitempty"`
	Width       string `url:"width,omitempty"`
	Height      string `url:"height,omitempty"`
	OrderId     string `url:"order_id,omitempty"`
	ItemId      string `url:"item_id,omitempty"`
	PromotionID string `url:"promotion_id,omitempty"`
	VideoId     string `url:"video_id,omitempty"`
}

// ReadSendMessageSession is a method on the ChatService struct. It sends a request to the server to send message to specific sessionID.
// opts: A pointer to a SendMessageParams struct containing the parameters for the SendMessage function.
// The function returns a pointer to an OpenSessionResposne struct containing the server's response, and an error, if there is one.
func (m *ChatService) SendMessage(ctx context.Context, opts *SendMessageParams) (res *SendMessageResponse, err error) {
	u, err := addOptions(ApiNames["SendMessage"], opts)
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

	json.Unmarshal([]byte(jsonData), &res)

	return res, nil
}

type GetListStickerResponse struct {
	Tab1 []Tab `json:"tab1"`
	Tab2 []Tab `json:"tab2"`
	Tab3 []Tab `json:"tab3"`
}
type Tab struct {
	Txt         string `json:"txt"`
	ImgURL      string `json:"imgUrl"`
	SmallImgURL string `json:"smallImgUrl"`
}

func (m *ChatService) GetListSticker() (*GetListStickerResponse, error) {
	jsonBytes, err := os.ReadFile("lazada/sticker-list.json")
	if err != nil {
		return nil, err
	}

	var response GetListStickerResponse
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
