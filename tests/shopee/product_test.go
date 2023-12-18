package tests

import (
	"fmt"
	"testing"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/shopee"
	"github.com/jarcoal/httpmock"
)

func Test_GetProduct(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/v2/product/get_item_base_info", app.APIURL),
		httpmock.NewBytesResponder(200, loadFixture("get_product_resp.json")))

	res, err := client.Product.GetProductById(123456, accessToken, shopee.GetProductParamRequest{
		ItemIDList:          []int{3400133011},
		NeedTaxInfo:         true,
		NeedComplaintPolicy: true,
	})

	if err != nil {
		t.Errorf("Product.GetProductResponse error: %s", err)
	}

	t.Logf("return tok: %#v", res)

	var expectedMsgID int64 = 3400133011
	if res.Response.ItemList[0].ItemID != expectedMsgID {
		t.Errorf("MessageList.MessageID returned %+v, expected %+v", res.Response.ItemList[0].ItemID, expectedMsgID)
	}
}

func Test_GetModelListt(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/v2/product/get_model_list", app.APIURL),
		httpmock.NewBytesResponder(200, loadFixture("get_model_list_resp.json")))

	res, err := client.Product.GetModelList(shopID, accessToken, 123)
	if err != nil {
		t.Errorf("Product.GetModelList error: %s", err)
	}

	t.Logf("Product.GetModelList: %#v", res)

	var expected uint64 = 2000458802
	if res.Response.Model[0].ModelID != expected {
		t.Errorf("ModelID returned %+v, expected %+v", res.Response.Model[0].ModelID, expected)
	}
}
