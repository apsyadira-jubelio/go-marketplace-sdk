package tokopedia

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

type TokopediaHTTPOpts struct {
	Token             string
	FsID              int64
	ShopID            int64
	SocksProxyAddress string
	APIURL            string
}

func NewTokopediaHTTPHandler(client *TokopediaClient, sockAddress string) (*TokopediaHTTPOpts, error) {
	if client == nil {
		return nil, errors.New("client is required")
	} else if client.appConfig.FsID == 0 {
		return nil, errors.New("fsID is required")
	}

	if sockAddress == "" {
		return nil, errors.New("sockAddress is required")
	}

	if client.appConfig.APIURL == "" {
		client.appConfig.APIURL = APIURL
	}

	intShopID, _ := strconv.Atoi(client.ShopID)
	return &TokopediaHTTPOpts{
		Token:             client.AccessToken,
		FsID:              int64(client.appConfig.FsID),
		ShopID:            int64(intShopID),
		SocksProxyAddress: fmt.Sprintf("socks5://%s", sockAddress),
		APIURL:            client.appConfig.APIURL,
	}, nil
}

func (opts *TokopediaHTTPOpts) GetListMessages(params GetMessagesParams) (*MessageResponse, error) {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 80 * time.Second,
	}
	proxyURL, err := url.Parse(opts.SocksProxyAddress)
	if err != nil {
		log.Println("error while parse socks address")
		return nil, err
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		log.Println("error while make transport dialer")
		return nil, err
	}
	client.Transport = &http.Transport{Dial: dialer.Dial}

	urlParam := fmt.Sprintf("%s/v1/chat/fs/%d/messages?page=%d&per_page=%d&shop_id=%d", opts.APIURL, opts.FsID, params.Page, params.PerPage, opts.ShopID)
	req, err := http.NewRequest(http.MethodGet, urlParam, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+opts.Token)
	req.Header.Add("Accept", "application/json")

	isDone := false
	maxRetries := 3

	var response MessageResponse
	for !isDone {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != http.StatusTooManyRequests {
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body, cause: %+v\n", err)
				return nil, err
			}
			log.Println("response:", map[string]interface{}{
				"body": string(respBytes),
				"code": resp.StatusCode,
			})
			return nil, errors.New("response not ok")
		}

		retryAfter := 0
		if resp.StatusCode == http.StatusTooManyRequests {
			headerReset := resp.Header.Get("X-Ratelimit-Full-Reset-After")
			f, _ := strconv.ParseFloat(headerReset, 64)
			retryAfter = int(f)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			maxRetries--
		} else {
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body, cause: %+v\n", err)
				return nil, err
			}

			if err := json.Unmarshal(respBytes, &response); err != nil {
				log.Printf("Error unmarshaling response, cause: %+v\n", err)
				log.Println("Response body:", string(respBytes))
				return nil, err
			}
			isDone = true
		}
		if maxRetries == 0 {
			isDone = true
		}
	}

	return &response, nil
}

func (opts *TokopediaHTTPOpts) GetProductInfo(params ProductParams) (*ProductInfoResponse, error) {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 80 * time.Second,
	}

	proxyURL, err := url.Parse(opts.SocksProxyAddress)
	if err != nil {
		log.Println("error while parse socks address")
		return nil, err
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		log.Println("error while dialer")
		return nil, err
	}
	client.Transport = &http.Transport{Dial: dialer.Dial}

	urlParam := fmt.Sprintf("%s/inventory/v1/fs/%d/product/info?product_id=%d", opts.APIURL, opts.FsID, params.ProductID)
	req, err := http.NewRequest(http.MethodGet, urlParam, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+opts.Token)
	req.Header.Add("Accept", "application/json")

	isDone := false
	maxRetries := 3

	var response ProductInfoResponse
	for !isDone {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != http.StatusTooManyRequests {
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body, cause: %+v\n", err)
				return nil, err
			}
			log.Println("response:", map[string]interface{}{
				"body": string(respBytes),
				"code": resp.StatusCode,
			})
			return nil, errors.New("response not ok")
		}

		retryAfter := 0
		if resp.StatusCode == http.StatusTooManyRequests {
			headerReset := resp.Header.Get("X-Ratelimit-Full-Reset-After")
			f, _ := strconv.ParseFloat(headerReset, 64)
			retryAfter = int(f)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			maxRetries--
		} else {
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body, cause: %+v\n", err)
				return nil, err
			}

			if err := json.Unmarshal(respBytes, &response); err != nil {
				log.Printf("Error unmarshaling response, cause: %+v\n", err)
				log.Println("Response body:", string(respBytes))
				return nil, err
			}
			isDone = true
		}
		if maxRetries == 0 {
			isDone = true
		}
	}

	return &response, nil
}

func (opts *TokopediaHTTPOpts) GetReplyTokopedia(params GetReplyListParams) (*ReplyListResponse, error) {
	var client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 80 * time.Second,
	}

	proxyURL, err := url.Parse(opts.SocksProxyAddress)
	if err != nil {
		log.Println("error while parse socks address")
		return nil, err
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		log.Println("error while parse direct")
		return nil, err
	}
	client.Transport = &http.Transport{Dial: dialer.Dial}

	urlParam := fmt.Sprintf("%s/v1/chat/fs/%d/messages/%d/replies?page=%d&per_page=%d&shop_id=%d", opts.APIURL, opts.FsID, params.MsgID, params.Page, params.PerPage, opts.ShopID)
	req, err := http.NewRequest(http.MethodGet, urlParam, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+opts.Token)
	req.Header.Add("Accept", "application/json")

	isDone := false
	maxRetries := 3

	var response ReplyListResponse
	for !isDone {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != http.StatusTooManyRequests {
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body, cause: %+v\n", err)
				return nil, err
			}
			log.Println("response:", map[string]interface{}{
				"body": string(respBytes),
				"code": resp.StatusCode,
			})
			return nil, errors.New("response not ok")
		}

		retryAfter := 0
		if resp.StatusCode == http.StatusTooManyRequests {
			headerReset := resp.Header.Get("X-Ratelimit-Full-Reset-After")
			f, _ := strconv.ParseFloat(headerReset, 64)
			retryAfter = int(f)
			time.Sleep(time.Duration(retryAfter) * time.Second)
			maxRetries--
		} else {
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body, cause: %+v\n", err)
				return nil, err
			}

			if err := json.Unmarshal(respBytes, &response); err != nil {
				log.Printf("Error unmarshaling response, cause: %+v\n", err)
				log.Println("Response body:", string(respBytes))
				return nil, err
			}
			isDone = true
		}
		if maxRetries == 0 {
			isDone = true
		}
	}

	return &response, nil
}
