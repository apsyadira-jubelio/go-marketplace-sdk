package shopee

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// Option is used to configure client with options
type Option func(c *ShopeeClient)

func WithRetry(retries int) Option {
	return func(c *ShopeeClient) {
		c.retries = retries
	}
}

func WithLogger(logger LeveledLoggerInterface) Option {
	return func(c *ShopeeClient) {
		c.log = logger
	}
}

func WithProxy(proxyHost string) Option {
	return func(c *ShopeeClient) {
		proxyURL, err := url.Parse(proxyHost)
		if err != nil {
			return
		}
		c.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
}

func WithTimeout(timeout int) Option {
	return func(c *ShopeeClient) {
		c.Client.Timeout = time.Duration(timeout) * time.Second
	}
}

func WithSocks5(socksAddress string) Option {
	return func(c *ShopeeClient) {
		proxyURL, err := url.Parse(fmt.Sprintf("socks5://%s", socksAddress))
		if err != nil {
			panic(err)
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			panic(err)
		}

		c.Client.Transport = &http.Transport{Dial: dialer.Dial}
	}
}
