package tiktok

import (
	"net/http"
	"net/url"
	"time"
)

// Option is used to configure client with options
type Option func(c *TiktokClient)

func WithLogger(logger LeveledLoggerInterface) Option {
	return func(c *TiktokClient) {
		c.log = logger
	}
}

func WithRetry(retries int) Option {
	return func(c *TiktokClient) {
		c.retries = retries
	}
}

func WithProxy(proxyHost string) Option {
	return func(c *TiktokClient) {
		proxyURL, err := url.Parse(proxyHost)
		if err != nil {
			return
		}
		c.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
}

func WithTimeout(timeout int) Option {
	return func(c *TiktokClient) {
		c.Client.Timeout = time.Duration(timeout) * time.Second
	}
}
