package dhmiddleware

import (
	"fmt"
	"github.com/cyongxue/magicbox/xhiris/xhid"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/kataras/iris/v12/context"
	"testing"
	"time"
)

func TestInitLimitUtil(t *testing.T) {
	logConf := xhlog.LoggerConf{
		Dir:        "E:\\GoglandProjects\\src\\github.com\\cyongxue\\magicbox\\xhiris\\xhlog\\logs",
		Prefix:     "test",
		Level:      "info",
		RotateSize: 1 * 1024 * 1024,
		Console:    true,
	}
	if err := xhlog.Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	irisCtx := context.NewContext(nil)
	irisCtx.Values().Set(xhid.TraceId, xhid.IdDriver(1))
	irisCtx.Values().Set(xhid.SpanId, xhid.MakeSpanId("hello world"))

	InitLimitUtil(10, 10, 10)

	unLimit := 0
	limit := 0
	for i := 0; i < 1000; i++ {
		if LimitEngine.CheckKeyLimit(irisCtx, "1234589") {
			limit++
		} else {
			unLimit++
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println(fmt.Sprintf("limit=%d; unLimit=%d", limit, unLimit))
	return
}

func BenchmarkInitLimitUtil(b *testing.B) {
	logConf := xhlog.LoggerConf{
		Dir:        "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhredis/logs",
		Prefix:     "test",
		Level:      "info",
		RotateSize: 1 * 1024 * 1024,
		Console:    true,
	}
	if err := xhlog.Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	irisCtx := context.NewContext(nil)
	irisCtx.Values().Set(xhid.TraceId, xhid.IdDriver(1))
	irisCtx.Values().Set(xhid.SpanId, xhid.MakeSpanId("hello world"))

	b.StopTimer()
	b.StartTimer()

	InitLimitUtil(10, 10, 10)

	unLimit := 0
	limit := 0
	for i := 0; i < b.N; i++ {
		if LimitEngine.CheckKeyLimit(irisCtx, "1234567890") {
			limit++
		} else {
			unLimit++
		}
	}

	fmt.Println(fmt.Sprintf("limit=%d; unLimit=%d", limit, unLimit))
	return
}
