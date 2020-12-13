package dhmiddleware

import (
	"bytes"
	"github.com/cyongxue/magicbox/xhiris/xhid"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

// DetectHoldUp 探测拦截处理
func DetectHoldUp() context.Handler {
	return func(ctx context.Context) {
		if ctx.Request().URL.String() == "/" {
			ctx.StatusCode(iris.StatusNotFound)
			return
		} else {
			ctx.Next()
		}
	}
}

// CorsPreMiddleware cors跨域资源共享处理
func CorsPreMiddleware() context.Handler {
	return func(ctx context.Context) {
		ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if ctx.Request().Method == iris.MethodOptions {
			ctx.StatusCode(204)
			return
		}
		ctx.Next()
	}
}

const (
	RequestStartTime = "_in_start_time"
)

// TraceSpanIdMiddleware trace id/span id/cspan id/start time处理
func TraceSpanIdMiddleware() context.Handler {
	return func(ctx context.Context) {
		// 提取trace id
		originTraceId := ctx.Request().Header.Get(xhid.TraceId)
		if originTraceId == "" {
			originTraceId = strconv.FormatInt(xhid.IdDriver(rand.Int63()).GetNextId(), 10)
		}
		ctx.Values().Set(xhid.TraceId, originTraceId)

		// 尝试提取parent id，如果没有，则说明是第一跳，则parent id为空
		parentId := ctx.Request().Header.Get(xhid.ParentId)
		if parentId != "" {
			ctx.Values().Set(xhid.ParentId, parentId)
		}

		// 尝试获取span id，如果没有，则自己生成一个span id
		spanId := ctx.Request().Header.Get(xhid.SpanId)
		if spanId == "" {
			spanId = xhid.MakeSpanId(ctx.Request().URL.String())
		}
		ctx.Values().Set(xhid.SpanId, spanId)

		// 记录下请求进来的时间
		ctx.Values().Set(RequestStartTime, time.Now())

		ctx.Next()
	}
}

// RequestInMiddleware 每一个请求的入口打印日志：_request_in
func RequestInMiddleware() context.Handler {
	return func(ctx context.Context) {
		var args interface{}
		switch ctx.Request().Method {
		case iris.MethodGet:
			args = ctx.URLParams()
		case iris.MethodPost:
			data, _ := ioutil.ReadAll(ctx.Request().Body)
			args = xhlog.TripSpaceAndReturn(data)
			// 要重新构建body
			ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}

		xhlog.Info(ctx, xhlog.OPRequestIn, map[string]interface{}{
			xhlog.Args: args,
		})

		ctx.Next()
	}
}
