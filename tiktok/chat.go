package tiktok

import (
	"fmt"
)

type ChatService interface {
	GetConversationMessages(params GetConversationMessagesParam, convesationID, shopChiper, shopID, accessToken string) (*GetConversationMessagesResponse, error)
	GetConversations(params GetConversationsParam, shopChiper, shopID, accessToken string) (*GetConversationsResponse, error)
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

func (s *ChatServiceOp) GetConversationMessages(params GetConversationMessagesParam, conversationID, shopChiper, shopID, accessToken string) (*GetConversationMessagesResponse, error) {
	path := fmt.Sprintf("/customer_service/%s/conversations/%s/messages", s.client.appConfig.Version, conversationID)
	resp := new(GetConversationMessagesResponse)
	err := s.client.WithShopID(shopID).WithShopChiper(shopChiper, s.client.appConfig.Version).WithAccessToken(accessToken).Get(path, resp, params)
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

func (s *ChatServiceOp) GetConversations(params GetConversationsParam, shopChiper, shopID, accessToken string) (*GetConversationsResponse, error) {
	path := fmt.Sprintf("/customer_service/%s/conversations", s.client.appConfig.Version)
	resp := new(GetConversationsResponse)
	err := s.client.WithShopID(shopID).WithShopChiper(shopChiper, s.client.appConfig.Version).WithAccessToken(accessToken).Get(path, resp, params)
	return resp, err
}
