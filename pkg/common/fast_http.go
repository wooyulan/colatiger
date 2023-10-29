package common

import (
	"crypto/tls"
	"github.com/valyala/fasthttp"
	"time"
)

func NewFastClient() *fasthttp.Client {
	return &fasthttp.Client{
		Name:                     "ColaTiger_Client",
		NoDefaultUserAgentHeader: true,
		TLSConfig:                &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost:          2000,
		MaxIdleConnDuration:      5 * time.Second,
		MaxConnDuration:          5 * time.Second,
		ReadTimeout:              5 * time.Second,
		WriteTimeout:             5 * time.Second,
		MaxConnWaitTimeout:       5 * time.Second,
	}
}
