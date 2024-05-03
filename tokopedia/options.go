package tokopedia

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// Option is used to configure client with options
type Option func(c *TokopediaClient)

func WithLogger(logger LeveledLoggerInterface) Option {
	return func(c *TokopediaClient) {
		c.log = logger
	}
}

func WithRetry(retries int) Option {
	return func(c *TokopediaClient) {
		c.retries = retries
	}
}

func WithProxy(proxyHost string) Option {
	return func(c *TokopediaClient) {
		proxyURL, err := url.Parse(proxyHost)
		if err != nil {
			return
		}
		c.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
}

func WithTimeout(timeout int) Option {
	return func(c *TokopediaClient) {
		c.Client.Timeout = time.Duration(timeout) * time.Second
	}
}

func WithSocks5(socksAddress string) Option {
	return func(c *TokopediaClient) {
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
