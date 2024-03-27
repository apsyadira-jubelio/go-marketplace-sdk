package tiktok

import (
	"fmt"
)

type ChatService interface {
	GetConversationMessages(conversationID string, param GetConversationMessagesParam) (*GetConversationMessagesResponse, error)
	GetConversations(params GetConversationsParam) (*GetConversationsResponse, error)
	SendMessageToConversationID(conversationID string, body SendMessageToConversationIDReq) (*SendMessageToConversationIDResp, error)
	ReadMessageConversationID(conversationID string) (*ReadMessageConversationIDResp, error)
	CreateConversation(body CreateConversationReq) (*CreateConversationResp, error)
	UploadBuyerMessagesImages(filename string) (*UploadMessagesImagesResp, error)
}

type ChatServiceOp struct {
	client *TiktokClient
}

type GetConversationMessagesParam struct {
	PageToken string `url:"page_token,omitempty"`
	PageSize  int    `url:"page_size"`
	Locale    string `url:"locale,omitempty"`
	SortOrder string `url:"sort_order,omitempty"`
	SortField string `url:"sort_field,omitempty"`
}

type GetConversationMessagesResponse struct {
	BaseResponse
	Data *DataConversationMessages `json:"data"`
}
type Sender struct {
	Avatar   string `json:"avatar"`
	ImUserID string `json:"im_user_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}
type MessagesConversation struct {
	Content    string  `json:"content"`
	CreateTime int     `json:"create_time"`
	ID         string  `json:"id"`
	IsVisible  bool    `json:"is_visible"`
	Sender     *Sender `json:"sender"`
	Type       string  `json:"type"`
}
type DataConversationMessages struct {
	Messages           []MessagesConversation `json:"messages"`
	NextPageToken      string                 `json:"next_page_token"`
	UnsupportedMsgTips string                 `json:"unsupported_msg_tips"`
}

func (s *ChatServiceOp) GetConversationMessages(conversationID string, params GetConversationMessagesParam) (*GetConversationMessagesResponse, error) {
	path := fmt.Sprintf("/customer_service/%s/conversations/%s/messages", s.client.appConfig.Version, conversationID)
	resp := new(GetConversationMessagesResponse)
	err := s.client.Get(path, resp, params)

	return resp, err
}

type GetConversationsParam struct {
	PageToken string `url:"page_token,omitempty"`
	PageSize  int    `url:"page_size"`
	Locale    string `url:"locale,omitempty"`
}

type GetConversationsResponse struct {
	BaseResponse
	Data *DataGetConversations `json:"data"`
}

type LatestMessage struct {
	Content    string  `json:"content"`
	CreateTime int     `json:"create_time"`
	ID         string  `json:"id"`
	IsVisible  bool    `json:"is_visible"`
	Sender     *Sender `json:"sender"`
	Type       string  `json:"type"`
}
type Participants struct {
	Avatar   string `json:"avatar"`
	ImUserID string `json:"im_user_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	UserID   string `json:"user_id"`
}
type Conversations struct {
	CanSendMessage   bool           `json:"can_send_message"`
	CreateTime       int            `json:"create_time"`
	ID               string         `json:"id"`
	LatestMessage    *LatestMessage `json:"latest_message"`
	ParticipantCount int            `json:"participant_count"`
	Participants     []Participants `json:"participants"`
	UnreadCount      int            `json:"unread_count"`
}
type DataGetConversations struct {
	Conversations []Conversations `json:"conversations"`
	NextPageToken string          `json:"next_page_token"`
}

func (s *ChatServiceOp) GetConversations(params GetConversationsParam) (*GetConversationsResponse, error) {
	path := fmt.Sprintf("/customer_service/%s/conversations", s.client.appConfig.Version)
	resp := new(GetConversationsResponse)
	err := s.client.
		Get(path, resp, params)

	return resp, err
}

const (
	TypeMessageText    = "TEXT"
	TypeMessageProduct = "PRODUCT_CARD"
	TypeMessageImage   = "IMAGE"
	TypeMessageOrder   = "ORDER_CARD"
)

type SendMessageToConversationIDReq struct {
	TypeMessage string `json:"type"`
	Content     string `json:"content"` // json stringfy
}

type SendMessageToConversationIDResp struct {
	BaseResponse
	Data *DataSendMsgResponse `json:"data"`
}

type DataSendMsgResponse struct {
	MessageID string `json:"message_id"`
}

func (s *ChatServiceOp) SendMessageToConversationID(conversationID string, body SendMessageToConversationIDReq) (*SendMessageToConversationIDResp, error) {
	path := fmt.Sprintf("/customer_service/%s/conversations/%s/messages", s.client.appConfig.Version, conversationID)

	resp := new(SendMessageToConversationIDResp)
	err := s.client.Post(path, body, resp)
	return resp, err
}

type ReadMessageConversationIDResp struct {
	BaseResponse
	Data interface{} `json:"data"`
}

func (s *ChatServiceOp) ReadMessageConversationID(conversationID string) (*ReadMessageConversationIDResp, error) {
	path := fmt.Sprintf("/customer_service/%s/conversations/%s/messages/read", s.client.appConfig.Version, conversationID)
	resp := new(ReadMessageConversationIDResp)
	err := s.client.Post(path, nil, resp)

	return resp, err
}

type UploadMessagesImagesResp struct {
	BaseResponse
	Data *DataUploadMessagesImages `json:"data"`
}

type DataUploadMessagesImages struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

func (s *ChatServiceOp) UploadBuyerMessagesImages(filename string) (*UploadMessagesImagesResp, error) {
	path := fmt.Sprintf("/customer_service/%s/images/upload", s.client.appConfig.Version)
	resp := new(UploadMessagesImagesResp)
	err := s.client.Upload(path, "data", filename, resp)
	return resp, err
}

type CreateConversationReq struct {
	BuyerUserID string `json:"buyer_user_id"`
}

type CreateConversationResp struct {
	BaseResponse
	Data *DataSendMsgResponse `json:"data"`
}

type DataCreateConversationResp struct {
	ConversationID string `json:"conversation_id"`
}

func (s *ChatServiceOp) CreateConversation(body CreateConversationReq) (*CreateConversationResp, error) {
	path := fmt.Sprintf("/customer_service/%s/conversations", s.client.appConfig.Version)
	resp := new(CreateConversationResp)
	err := s.client.Post(path, body, resp)
	return resp, err
}
