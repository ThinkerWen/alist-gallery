package net

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"sync"
	"time"
)

var (
	GlobalClient *resty.Client
	once         sync.Once
)

func init() {
	once.Do(func() {
		GlobalClient = resty.New().
			SetRetryCount(3).
			SetRetryWaitTime(500 * time.Millisecond).
			SetRetryMaxWaitTime(2 * time.Second).
			SetTimeout(30 * time.Second).
			SetTransport(&http.Transport{
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			})
	})
}
