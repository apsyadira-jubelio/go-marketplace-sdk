package tokopedia

import (
	"fmt"
)

type ProductInfoResponse struct {
	BaseResponse
	Data []ProductData `json:"data"`
}

type ProductParams struct {
	ProductID int `url:"product_id"`
}

type ProductData struct {
	Basic struct {
		ProductID       int    `json:"productID"`
		ShopID          int    `json:"shopID"`
		Status          int    `json:"status"`
		Name            string `json:"Name"`
		Condition       int    `json:"condition"`
		ChildCategoryID int    `json:"childCategoryID"`
		ShortDesc       string `json:"shortDesc"`
	} `json:"basic"`
	Price struct {
		Value          int `json:"value"`
		Currency       int `json:"currency"`
		LastUpdateUnix int `json:"LastUpdateUnix"`
		Idr            int `json:"idr"`
	} `json:"price"`
	Weight struct {
		Value int `json:"value"`
		Unit  int `json:"unit"`
	} `json:"weight"`
	Stock struct {
		Value        int    `json:"value"`
		StockWording string `json:"stockWording"`
	} `json:"stock"`
	MainStock    int `json:"main_stock"`
	ReserveStock int `json:"reserve_stock"`
	Variant      struct {
		IsParent   bool  `json:"isParent"`
		IsVariant  bool  `json:"isVariant"`
		ChildrenID []int `json:"childrenID"`
	} `json:"variant"`
	Menu struct {
		ID   int    `json:"id"`
		Name string `json:"Name"`
	} `json:"menu"`
	ExtraAttribute struct {
		MinOrder           int  `json:"minOrder"`
		LastUpdateCategory int  `json:"lastUpdateCategory"`
		IsEligibleCOD      bool `json:"isEligibleCOD"`
	} `json:"extraAttribute"`
	CategoryTree []struct {
		ID            int    `json:"id"`
		Name          string `json:"Name"`
		Title         string `json:"title"`
		BreadcrumbURL string `json:"breadcrumbURL"`
	} `json:"categoryTree"`
	Pictures []struct {
		PicID        int    `json:"picID"`
		FileName     string `json:"fileName"`
		FilePath     string `json:"filePath"`
		Status       int    `json:"status"`
		OriginalURL  string `json:"OriginalURL"`
		ThumbnailURL string `json:"ThumbnailURL"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		URL300       string `json:"URL300"`
	} `json:"pictures"`
	GMStats struct {
		TransactionSuccess int `json:"transactionSuccess"`
		TransactionReject  int `json:"transactionReject"`
		CountSold          int `json:"countSold"`
	} `json:"GMStats"`
	Stats struct {
		CountView int `json:"countView"`
	} `json:"stats"`
	Other struct {
		Sku       string `json:"sku"`
		URL       string `json:"url"`
		MobileURL string `json:"mobileURL"`
	}
}

type ProductService interface {
	GetProductInfo(token string, productID int) (res *ProductInfoResponse, err error)
}

type ProductServiceOp struct {
	client *TokopediaClient
}

func (p *ProductServiceOp) GetProductInfo(token string, productID int) (res *ProductInfoResponse, err error) {
	path := fmt.Sprintf("/inventory/v1/fs/%d/product/info", p.client.appConfig.FsID)

	resp := new(ProductInfoResponse)

	if productID == 0 {
		return nil, fmt.Errorf("product_id is required")
	}

	params := ProductParams{
		ProductID: productID,
	}

	err = p.client.WithAccessToken(token).Get(path, resp, params)
	return resp, err

}
