package tiktok

import "fmt"

type ProductService interface {
	GetProductInfo(productID string) (*GetProductInfoResponse, error)
}

type ProductServiceOp struct {
	client *TiktokClient
}

type GetProductParamRequest struct {
	ProductID string `url:"product_id"`
}

type GetProductInfoResponse struct {
	BaseResponse
	Data *ProductData `json:"data"`
}

type ProductData struct {
	Brand             Brand              `json:"brand"`
	CategoryChains    []CategoryChain    `json:"category_chains"`
	CreateTime        int64              `json:"create_time"`
	Description       string             `json:"description"`
	ID                string             `json:"id"`
	IsCodAllowed      bool               `json:"is_cod_allowed"`
	MainImages        []MainImage        `json:"main_images"`
	PackageDimensions *PackageDimensions `json:"package_dimensions"`
	PackageWeight     *PackageWeight     `json:"package_weight"`
	ProductAttributes []ProductAttribute `json:"product_attributes"`
	Skus              []Skus             `json:"skus"`
	Status            string             `json:"status"`
	Title             string             `json:"title"`
	UpdateTime        int64              `json:"update_time"`
}

type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryChain struct {
	ID        string `json:"id"`
	IsLeaf    bool   `json:"is_leaf"`
	LocalName string `json:"local_name"`
	ParentID  string `json:"parent_id"`
}

type MainImage struct {
	Height    int64    `json:"height"`
	ThumbUrls []string `json:"thumb_urls"`
	URI       string   `json:"uri"`
	Urls      []string `json:"urls"`
	Width     int64    `json:"width"`
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

type ProductAttribute struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Values []Brand `json:"values"`
}

type Skus struct {
	ID              string         `json:"id"`
	IdentifierCode  IdentifierCode `json:"identifier_code"`
	Inventory       []Inventory    `json:"inventory"`
	Price           *Price         `json:"price"`
	SalesAttributes []interface{}  `json:"sales_attributes"`
	SellerSku       string         `json:"seller_sku"`
}

type IdentifierCode struct {
	Type string `json:"type"`
}

type Inventory struct {
	Quantity    int64  `json:"quantity"`
	WarehouseID string `json:"warehouse_id"`
}

type Price struct {
	Currency          string `json:"currency"`
	SalePrice         string `json:"sale_price"`
	TaxExclusivePrice string `json:"tax_exclusive_price"`
}

func (p *ProductServiceOp) GetProductInfo(productID string) (*GetProductInfoResponse, error) {
	path := fmt.Sprintf("/product/%s/products/%s", p.client.appConfig.Version, productID)

	resp := new(GetProductInfoResponse)
	err := p.client.Get(path, resp, nil)

	return resp, err
}
