package goclient

import (
	"net/http"
	"time"
)

type Config struct {
	BaseURL               string
	Timeout               time.Duration
	GlobalHeaders         map[string]string
	Interceptor           http.RoundTripper
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	DisableKeepAlives     bool
	DisableCompression    bool
	ResponseHeaderTimeout time.Duration
}

type Option func(*Config)

func defaultConfig(config ...Config) Config {
	cfg := Config{
		Timeout:               30 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	if len(config) > 0 {
		cfg = config[0]
	}

	return cfg
}

func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithGlobalHeaders(headers map[string]string) Option {
	return func(c *Config) {
		c.GlobalHeaders = headers
	}
}

func WithMaxIdleConns(n int) Option {
	return func(c *Config) {
		c.MaxIdleConns = n
	}
}

func WithMaxIdleConnsPerHost(n int) Option {
	return func(c *Config) {
		c.MaxIdleConnsPerHost = n
	}
}

func WithDisableKeepAlives(disable bool) Option {
	return func(c *Config) {
		c.DisableKeepAlives = disable
	}
}

func WithDisableCompression(disable bool) Option {
	return func(c *Config) {
		c.DisableCompression = disable
	}
}
