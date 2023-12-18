package tokopedia

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type ResponseHeader struct {
	ProcessTime int    `json:"process_time"`
	Message     string `json:"message"`
	Reason      string `json:"reason"`
	ErrorCode   string `json:"error_code"`
}

type TokopediaAuthResponse struct {
	AccessToken   string `json:"access_token"`
	EventCode     string `json:"event_code"`
	ExpiresIn     int64  `json:"expires_in"`
	LastLoginType string `json:"last_login_type"`
	SqCheck       bool   `json:"sq_check"`
	TokenType     string `json:"token_type"`
}

type TokopediaMessageResponse struct {
	ResponseHeader
	Data TokoepdiaMessageResponseData `json:"data"`
}

type TokoepdiaMessageResponseData struct {
	MsgID      int64           `json:"msg_id"`
	SenderID   int64           `json:"sender_id"`
	Role       int             `json:"role"`
	Msg        string          `json:"msg"`
	ReplyTime  int64           `json:"reply_time"`
	From       string          `json:"from"`
	Attachemnt json.RawMessage `json:"attachment"`
}

type TokopediaMessageText struct {
	ShopId  int64  `json:"shop_id"`
	Message string `json:"message"`
}

type SendReplyAttachment struct {
	ShopId         int64 `json:"shop_id"`
	AttachmentType int   `json:"attachment_type"`
	Payload        struct {
		Thumbnail  string `json:"thumbnail"`
		Identifier string `json:"identifier"`
		Title      string `json:"title"`
		Price      string `json:"price"`
		Url        string `json:"url"`
	} `json:"payload"`
}

// type TokopediaUsecase interface {
// 	SendMessageText(c *fiber.Ctx, fsId, msgId int64, token string, data TokopediaMessageText) error
// 	IntegrateAccount(c *fiber.Ctx, payload ChannelData) error
// 	SaveToken(c *fiber.Ctx, payload ChannelData, tokopediaToken TokopediaAuthResponse) error
// }

// TokopediaAPI represents a struct to interact with the Tokopedia API.
// It contains a resty client to make the HTTP requests, a token for authorization,
// and a fsID that could represent some form of identification or status.
type TokopediaClient struct {
	Client *resty.Client
	token  string
	fsID   *int64

	Auth AuthService
	Chat ChatService
}

// NewTokopediaApi is a constructor function for creating a new instance of TokopediaAPI.
// It accepts three parameters: isAuth (a boolean to indicate if authentication is required),
// token (a string representing the authorization token), and fsID (a pointer to an int64 which could be an identifier).
// The function returns a pointer to a newly created TokopediaAPI instance.
func NewClient(isAuth bool, token string, fsID *int64) *TokopediaClient {

	// Initialize a new resty client using the TokopediaClient function
	rc := HTTPTokopediaApi(isAuth)

	// If a token is provided, add it to the Authorization header of the resty client
	if token != "" {
		rc.SetHeader("Authorization", "Bearer "+token)
	}

	c := &TokopediaClient{
		Client: rc,
		token:  token,
		fsID:   fsID,
	}

	c.Auth = &AuthServiceOp{client: c}
	c.Chat = &ChatServiceOp{client: c}

	// Return a new TokopediaAPI instance with the client, token, and fsID fields set
	return c
}
