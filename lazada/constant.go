package lazada

// API Names are all the paths to the various API calls that we use
var ApiNames = map[string]string{
	"AccessToken":        "https://auth.lazada.com/rest/auth/token/create",
	"RefreshToken":       "https://auth.lazada.com/rest/auth/token/refresh",
	"GetBrands":          "/brands/get",
	"CategoryTree":       "/category/tree/get",
	"ImageMigrate":       "/image/migrate",
	"CategoryAttributes": "/category/attributes/get",
	"CreateProduct":      "/product/create",
	"UpdateProduct":      "/product/update",
	"GetProducts":        "/products/get",
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
