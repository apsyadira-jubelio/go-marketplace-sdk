package main

import (
	"log"
	"os"

	"github.com/apsyadira-jubelio/go-marketplace-sdk/tokopedia"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)

	app := tokopedia.AppConfig{
		FsID:   14332,
		APIURL: "https://fs.tokopedia.net",
	}

	withProxy := tokopedia.WithSocks5(os.Getenv("SOCKS_ADDRESS"))
	client := tokopedia.NewClient(app, withProxy)

	//c:QPVk0IZVQp-Ie2zSgSf1cg
	// 2545808044
	resp, err := client.Product.GetProductInfo("c:Q07-p3k8SvKrWlEMoE_wUg", 0)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(spew.Sdump(resp))
}
