package lazada

// API Names are all the paths to the various API calls that we use
var ApiNames = map[string]string{
	"AccessToken":            "https://auth.lazada.com/rest/auth/token/create",
	"RefreshToken":           "https://auth.lazada.com/rest/auth/token/refresh",
	"GetBrands":              "/brands/get",
	"CategoryTree":           "/category/tree/get",
	"ImageMigrate":           "/image/migrate",
	"CategoryAttributes":     "/category/attributes/get",
	"CreateProduct":          "/product/create",
	"UpdateProduct":          "/product/update",
	"GetProducts":            "/products/get",
	"GetSessionList":         "/im/session/list",
	"GetSessionDetail":       "/im/session/get",
	"GetMessageList":         "/im/message/list",
	"MessageRecall":          "/im/message/recall",
	"OpenSession":            "/im/session/open",
	"ReadSession":            "/im/session/read",
	"SendMessage":            "/im/message/send",
	"GetMultipleOrdersItems": "/orders/items/get",
	"GetOrders":              "/orders/get",
	"GetVideo":               "/media/video/get",
	"GetVouchers":            "/promotion/vouchers/get",
	"OrderTrace":             "/logistic/order/trace",
}

type Region string

const (
	Philippines = "ph"
	Bangladesh  = "bd"
	Thailand    = "th"
	Vietnam     = "vn"
	Pakistan    = "pk"
	Singapore   = "sg"
	Nepal       = "np"
	Indonesia   = "id"
	Myanmar     = "mm"
	Malaysia    = "my"
)

// endpoints maps a regions shortcode to its URL
var endpoints = map[Region]string{
	Philippines: "https://api.lazada.com.ph/",
	Bangladesh:  "https://api.daraz.com.bd/",
	Thailand:    "https://api.lazada.co.th/",
	Vietnam:     "https://api.lazada.vn/",
	Pakistan:    "https://api.daraz.pk/",
	Singapore:   "https://api.lazada.sg/",
	Nepal:       "https://api.daraz.com.np/",
	Indonesia:   "https://api.lazada.co.id/",
	Myanmar:     "https://api.shop.com.mm/",
	Malaysia:    "https://api.lazada.com.my/",
}

// Message templpate ID to send message
const (
	NormalTextMessage         = 1
	PictureMessage            = 3
	EmojiMessage              = 4
	ItemMessage               = 10006
	OrderMessage              = 10007
	VoucherMessage            = 10008
	InviteBuyersFollowMessage = 10010
	VideoMessage              = 6
)
