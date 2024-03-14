package shopee

type ProductService interface {
	GetProductById(shopID uint64, token string, params GetProductParamRequest) (*GetProductResponse, error)
	GetModelList(shopID uint64, token string, itemID uint64) (*GetModelListResponse, error)
	GetProductWithSearch(shopID uint64, token string, paramRequest GetProductWithSearchRequest) (*GetProductWithSearchResponse, error)
	GetProductlList(shopID uint64, token string, paramRequest GetProductListParamRequest) (*GetProductListResponse, error)
}

type GetProductResponse struct {
	BaseResponse
	Response ItemListResponse `json:"response"`
}

type ItemListResponse struct {
	ItemList []ItemListData `json:"item_list"`
}

type ItemListData struct {
	ItemID          int64           `json:"item_id"`
	CategoryID      int64           `json:"category_id"`
	ItemName        string          `json:"item_name"`
	ItemSku         string          `json:"item_sku"`
	CreateTime      int64           `json:"create_time"`
	UpdateTime      int64           `json:"update_time"`
	AttributeList   []AttributeList `json:"attribute_list"`
	PriceInfo       []PriceInfo     `json:"price_info"`
	StockInfoV2     StockInfoV2     `json:"stock_info_v2"`
	Image           Image           `json:"image"`
	Weight          string          `json:"weight"`
	Dimension       Dimension       `json:"dimension"`
	LogisticInfo    []LogisticInfo  `json:"logistic_info"`
	PreOrder        PreOrder        `json:"pre_order"`
	Condition       string          `json:"condition"`
	SizeChart       string          `json:"size_chart"`
	ItemStatus      string          `json:"item_status"`
	HasModel        bool            `json:"has_model"`
	PromotionID     int64           `json:"promotion_id"`
	Brand           Brand           `json:"brand"`
	TaxInfo         TaxInfo         `json:"tax_info"`
	DescriptionType string          `json:"description_type"`
	DescriptionInfo DescriptionInfo `json:"description_info"`
}

type AttributeList struct {
	AttributeID           int64                `json:"attribute_id"`
	OriginalAttributeName string               `json:"original_attribute_name"`
	IsMandatory           bool                 `json:"is_mandatory"`
	AttributeValueList    []AttributeValueList `json:"attribute_value_list"`
}

type AttributeValueList struct {
	ValueID           int64  `json:"value_id"`
	OriginalValueName string `json:"original_value_name"`
	ValueUnit         string `json:"value_unit"`
}

type Brand struct {
	BrandID           int64  `json:"brand_id"`
	OriginalBrandName string `json:"original_brand_name"`
}

type DescriptionInfo struct {
	ExtendedDescription ExtendedDescription `json:"extended_description"`
}

type ExtendedDescription struct {
	FieldList []FieldList `json:"field_list"`
}

type FieldList struct {
	FieldType string     `json:"field_type"`
	Text      *string    `json:"text,omitempty"`
	ImageInfo *ImageInfo `json:"image_info,omitempty"`
}

type ImageInfo struct {
	ImageID  string `json:"image_id"`
	ImageURL string `json:"image_url"`
}

type Dimension struct {
	PackageLength int64 `json:"package_length"`
	PackageWidth  int64 `json:"package_width"`
	PackageHeight int64 `json:"package_height"`
}

type Image struct {
	ImageURLList []string `json:"image_url_list"`
	ImageIDList  []string `json:"image_id_list"`
}

type LogisticInfo struct {
	LogisticID           int64    `json:"logistic_id"`
	LogisticName         string   `json:"logistic_name"`
	Enabled              bool     `json:"enabled"`
	ShippingFee          *int64   `json:"shipping_fee,omitempty"`
	IsFree               bool     `json:"is_free"`
	EstimatedShippingFee *float64 `json:"estimated_shipping_fee,omitempty"`
}

type PreOrder struct {
	IsPreOrder bool  `json:"is_pre_order"`
	DaysToShip int64 `json:"days_to_ship"`
}

type PriceInfo struct {
	Currency                     string  `json:"currency"`
	OriginalPrice                float64 `json:"original_price"`
	CurrentPrice                 float64 `json:"current_price"`
	InflatedPriceOfOriginalPrice float64 `json:"inflated_price_of_original_price"`
	InflatedPriceOfCurrentPrice  float64 `json:"inflated_price_of_current_price"`
	SipItemPrice                 float64 `json:"sip_item_price"`
	SipItemPriceSource           string  `json:"sip_item_price_source"`
}

type StockInfoV2 struct {
	SummaryInfo SummaryInfo   `json:"summary_info"`
	SellerStock []SellerStock `json:"seller_stock"`
}

type SellerStock struct {
	LocationID string `json:"location_id"`
	Stock      int64  `json:"stock"`
}

type SummaryInfo struct {
	TotalReservedStock  int64 `json:"total_reserved_stock"`
	TotalAvailableStock int64 `json:"total_available_stock"`
}

type TaxInfo struct {
	Ncm           int64 `json:"ncm"`
	SameStateCfop int64 `json:"same_state_cfop"`
	DiffStateCfop int64 `json:"diff_state_cfop"`
	Csosn         int64 `json:"csosn"`
	Origin        int64 `json:"origin"`
}

type GetProductParamRequest struct {
	ItemIDList          []int `url:"item_id_list"`
	NeedTaxInfo         bool  `url:"need_tax_info"`
	NeedComplaintPolicy bool  `url:"need_complaint_policy"`
}

type ProductServiceOp struct {
	client *ShopeeClient
}

func (s *ProductServiceOp) GetProductById(shopID uint64, token string, params GetProductParamRequest) (*GetProductResponse, error) {
	path := "/product/get_item_base_info"

	resp := new(GetProductResponse)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, params)
	return resp, err
}

type GetModelListRequest struct {
	ItemID uint64 `url:"item_id"`
}

type GetModelListResponse struct {
	BaseResponse

	Response GetModelListResponseData `json:"response"`
}

type TierVariation struct {
	Name       string                `json:"name"`
	OptionList []TierVariationOption `json:"option_list"`
}

type TierVariationOption struct {
	Option string                    `json:"option"`
	Image  *TierVariationOptionImage `json:"image,omitempty"`
}

type TierVariationOptionImage struct {
	ImageID  string `json:"image_id"`
	ImageURL string `json:"image_url"`
}

type GetModelListResponseData struct {
	TierVariation []TierVariation `json:"tier_variation"`
	Model         []Model         `json:"model"`
}
type Model struct {
	TierIndex   []int       `json:"tier_index"`
	ModelID     uint64      `json:"model_id"`
	ModelSKU    string      `json:"model_sku"`
	StockInfo   []StockInfo `json:"stock_info"`
	PriceInfo   []PriceInfo `json:"price_info"`
	PromotionID uint64      `json:"promotion_id"`
}

type StockInfo struct {
	StockType       int    `json:"stock_type"`
	StockLocationID string `json:"stock_location_id"`
	NormalStock     int    `json:"normal_stock"`
	CurrentStock    int    `json:"current_stock"`
	ReservedStock   int    `json:"reserved_stock"`
}

func (s *ProductServiceOp) GetModelList(shopID uint64, token string, itemID uint64) (*GetModelListResponse, error) {
	path := "/product/get_model_list"

	opt := GetModelListRequest{
		ItemID: itemID,
	}

	resp := new(GetModelListResponse)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, opt)
	return resp, err
}

type GetProductListResponse struct {
	BaseResponse

	Response ProductListData `json:"response"`
}

type GetProductListParamRequest struct {
	Offset     int    `url:"offset"`
	PageSize   int    `url:"page_size"`
	ItemStatus string `url:"item_status"`
}

type ProductListData struct {
	Item        []ItemProductList `json:"item"`
	TotalCount  int               `json:"total_count"`
	HasNextPage bool              `json:"has_next_page"`
	NextOffset  int               `json:"next_offset"`
}

type ItemProductList struct {
	ItemID     int64  `json:"item_id"`
	ItemStatus string `json:"item_status"`
	UpdateTime int    `json:"update_time"`
}

func (s *ProductServiceOp) GetProductlList(shopID uint64, token string, paramRequest GetProductListParamRequest) (*GetProductListResponse, error) {
	path := "/product/get_item_list"

	resp := new(GetProductListResponse)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, paramRequest)
	return resp, err
}

type GetProductWithSearchResponse struct {
	BaseResponse

	Response ProductListData `json:"response"`
}

type GetProductWithSearchRequest struct {
	Offset     int    `url:"offset"`
	PageSize   int    `url:"page_size"`
	ItemStatus string `url:"item_status"`
	ItemName   string `url:"item_name"`
}

type GetProductWithSearchData struct {
	Item        []ItemProductListWithSearch `json:"item"`
	TotalCount  int                         `json:"total_count"`
	HasNextPage bool                        `json:"has_next_page"`
	NextOffset  int                         `json:"next_offset"`
}

type ItemProductListWithSearch struct {
	ItemID     int64  `json:"item_id"`
	ItemStatus string `json:"item_status"`
	UpdateTime int    `json:"update_time"`
}

func (s *ProductServiceOp) GetProductWithSearch(shopID uint64, token string, paramRequest GetProductWithSearchRequest) (*GetProductWithSearchResponse, error) {
	path := "/product/search_item"

	resp := new(GetProductWithSearchResponse)
	err := s.client.WithShop(uint64(shopID), token).Get(path, resp, paramRequest)
	return resp, err
}
