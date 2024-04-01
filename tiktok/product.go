package tiktok

import "fmt"

type ProductService interface {
	GetProductByID(conversationID string) (*GetProductByIDResp, error)
}

type ProductServiceOp struct {
	client *TiktokClient
}

type GetProductByIDResp struct {
	BaseResponse
	Data DataProduct `json:"data"`
}

type DataProduct struct {
	AuditFailedReasons []AuditFailedReasons `json:"audit_failed_reasons"`
	CategoryChains     []CategoryChains     `json:"category_chains"`
	CreateTime         int                  `json:"create_time"`
	Description        string               `json:"description"`
	ID                 string               `json:"id"`
	IsCodAllowed       bool                 `json:"is_cod_allowed"`
	MainImages         []MainImages         `json:"main_images"`
	PackageDimensions  *PackageDimensions   `json:"package_dimensions"`
	PackageWeight      *PackageWeight       `json:"package_weight"`
	Skus               []Skus               `json:"skus"`
	Status             string               `json:"status"`
	Title              string               `json:"title"`
	UpdateTime         int                  `json:"update_time"`
}

type AuditFailedReasons struct {
	Position    string   `json:"position"`
	Reasons     []string `json:"reasons"`
	Suggestions []string `json:"suggestions"`
}
type CategoryChains struct {
	ID        string `json:"id"`
	IsLeaf    bool   `json:"is_leaf"`
	LocalName string `json:"local_name"`
	ParentID  string `json:"parent_id"`
}

type MainImages struct {
	Height    int      `json:"height"`
	ThumbUrls []string `json:"thumb_urls"`
	URI       string   `json:"uri"`
	Urls      []string `json:"urls"`
	Width     int      `json:"width"`
}

type PackageDimensions struct {
	Height string `json:"height"`
	Length string `json:"length"`
	Unit   string `json:"unit"`
	Width  string `json:"width"`
}

type PackageWeight struct {
	Unit  string `json:"unit"`
	Value string `json:"value"`
}
type Inventory struct {
	Quantity    int    `json:"quantity"`
	WarehouseID string `json:"warehouse_id"`
}

type Price struct {
	Currency          string `json:"currency"`
	SalePrice         string `json:"sale_price"`
	TaxExclusivePrice string `json:"tax_exclusive_price"`
}

type SkuImg struct {
	Height    int      `json:"height"`
	ThumbUrls []string `json:"thumb_urls"`
	URI       string   `json:"uri"`
	Urls      []string `json:"urls"`
	Width     int      `json:"width"`
}

type SalesAttributes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SkuImg    SkuImg `json:"sku_img,omitempty"`
	ValueID   string `json:"value_id"`
	ValueName string `json:"value_name"`
}

type Skus struct {
	ID              string            `json:"id"`
	Inventory       []Inventory       `json:"inventory"`
	Price           *Price            `json:"price"`
	SalesAttributes []SalesAttributes `json:"sales_attributes"`
	SellerSku       string            `json:"seller_sku"`
}

func (p *ProductServiceOp) GetProductByID(productID string) (*GetProductByIDResp, error) {
	path := fmt.Sprintf("/product/%s/products/%s", p.client.appConfig.Version, productID)
	resp := new(GetProductByIDResp)
	err := p.client.Get(path, resp, nil)
	return resp, err
}
