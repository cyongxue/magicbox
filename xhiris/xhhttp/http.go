package xhhttp

import (
	"github.com/kataras/iris/v12"
	"net/http"
	"time"
)

type Http interface {
	Send(ctx iris.Context, reqBody []byte, head *http.Header) (int, []byte)
	SendWithTimeout(ctx iris.Context, reqBody []byte, head *http.Header, timeoutSec time.Duration) (int, []byte)
	HttpWithOption(ctx iris.Context, reqBody []byte, head *http.Header, timeoutSec time.Duration) (int, []byte)
}

type Method string

const (
	POST Method = "POST"
	GET         = "GET"
)
