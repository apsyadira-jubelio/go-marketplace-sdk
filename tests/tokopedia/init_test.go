package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

var (
	client *resty.Client
)

func setup() {
	client = resty.New()

	httpmock.ActivateNonDefault(client.GetClient())
}

func teardown() {
	httpmock.DeactivateAndReset()
}

func loadFixture(filename string) []byte {
	f, err := ioutil.ReadFile("../../mockdata/tokopedia/" + filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot load fixture %v", filename))
	}

	return f
}

func loadMockData(filename string, out interface{}) {
	f, err := ioutil.ReadFile("../../mockdata/tokopedia/" + filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot load fixture %v", filename))
	}
	if err := json.Unmarshal(f, &out); err != nil {
		panic(fmt.Sprintf("decode mock data error: %s", err))
	}
}
