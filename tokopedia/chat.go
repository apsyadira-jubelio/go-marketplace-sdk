package tokopedia

import (
	"fmt"
)

type MessageResponse struct {
	BaseResponse
	Data []MessageData `json:"data"`
}

type MessageData struct {
	MessageKey string `json:"message_key"`
	MsgID      int    `json:"msg_id"`
	Attributes struct {
		Contact struct {
			ID         int    `json:"id"`
			Role       string `json:"role"`
			Attributes struct {
				Name      string `json:"Name"`
				Tag       string `json:"tag"`
				Thumbnail string `json:"thumbnail"`
			} `json:"attributes"`
		} `json:"contact"`
		LastReplyMsg  string `json:"last_reply_msg"`
		LastReplyTime int64  `json:"last_reply_time"`
		ReadStatus    int    `json:"read_status"`
		Unreads       int    `json:"unreads"`
		PinStatus     int    `json:"pin_status"`
	} `json:"attributes"`
}

type GetMessagesParams struct {
	Page    int   `url:"page,required"`
	PerPage int   `url:"per_page,required"`
	FsID    int64 `url:"fs_id,required"`
	ShopID  int   `url:"shop_id,required"`
}

type GetReplyListParams struct {
	ShopID  int `url:"shop_id,required"`
	MsgID   int `url:"msg_id,required"`
	Page    int `url:"page,required"`
	PerPage int `url:"per_page,required"`
}

type ReplyListResponse struct {
	BaseResponse
	Data []ReplyData `json:"data"`
}

type ReplyData struct {
	MsgID         int    `json:"msg_id"`
	SenderID      int    `json:"sender_id"`
	Role          string `json:"role"`
	Msg           string `json:"msg"`
	ReplyTime     int64  `json:"reply_time"`
	ReplyID       int    `json:"reply_id"`
	SenderName    string `json:"sender_name"`
	ReadStatus    int    `json:"read_status"`
	ReadTime      int64  `json:"read_time"`
	Status        int    `json:"status"`
	AttachmentID  int    `json:"attachment_id"`
	MessageIsRead bool   `json:"message_is_read"`
	IsOpposite    bool   `json:"is_opposite"`
	IsFirstReply  bool   `json:"is_first_reply"`
	IsReported    bool   `json:"is_reported"`
}

type SendMessageBody struct {
	Message        string `json:"message"`
	MsgID          int    `json:"msg_id"`
	ShopID         int    `json:"shop_id"`
	AttachmentType int    `json:"attachment_type"`
}

type SendMessageResponse struct {
	BaseResponse
	Data SendMessageResponseData `json:"data"`
}

type SendMessageResponseData struct {
	MsgID      int64  `json:"msg_id"`
	SenderID   int    `json:"sender_id"`
	Role       int    `json:"role"`
	Msg        string `json:"msg"`
	ReplyTime  int64  `json:"reply_time"`
	From       string `json:"from"`
	Attachment struct {
	} `json:"attachment"`
}

type ChatService interface {
	GetMessagesList(token string, params GetMessagesParams) (res *MessageResponse, err error)
	GetReplyList(token string, params GetReplyListParams) (res *ReplyListResponse, err error)
	SendMessage(token string, msgID int, body SendMessageBody) (res *SendMessageResponse, err error)
}

type ChatServiceOp struct {
	client *TokopediaClient
}

func (s *ChatServiceOp) GetMessagesList(token string, params GetMessagesParams) (res *MessageResponse, err error) {

	path := fmt.Sprintf("/v1/chat/fs/%d/messages", s.client.appConfig.FsID)

	if s.client.appConfig.FsID == 0 {
		return nil, fmt.Errorf("fs_id is required")
	}

	resp := new(MessageResponse)
	err = s.client.WithAccessToken(token).Get(path, resp, params)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func (s *ChatServiceOp) GetReplyList(token string, params GetReplyListParams) (res *ReplyListResponse, err error) {

	path := fmt.Sprintf("/v1/chat/fs/%d/messages/%d/replies", s.client.appConfig.FsID, params.MsgID)
	resp := new(ReplyListResponse)
	err = s.client.WithAccessToken(token).Get(path, resp, params)
	return resp, err

}

func (s *ChatServiceOp) SendMessage(token string, msgID int, body SendMessageBody) (res *SendMessageResponse, err error) {
	path := fmt.Sprintf("/v1/chat/fs/%d/messages/%d/reply", s.client.appConfig.FsID, msgID)
	resp := new(SendMessageResponse)
	err = s.client.WithAccessToken(token).Post(path, body, resp)
	return resp, err

}
