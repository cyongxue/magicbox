package xhhttp

import (
	"crypto/tls"
	"fmt"
	"github.com/cyongxue/magicbox/xhiris/xhid"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type NetHttp struct {
	Method  Method
	Url     string
	IsHttps bool
}

// Send 发送http请求
func (n *NetHttp) Send(ctx iris.Context, reqBody []byte, head *http.Header) (int, []byte) {
	startTime := time.Now()
	timestamp := xhlog.NowUSec()

	transport := http.Transport{
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * 10, //超时时间从5s调整到8s，获取设备树接口在组织下设备很多时易超时。
	}
	httpClient := http.Client{
		Timeout:   10 * time.Second,
		Transport: &transport,
	}

	// request 准备
	req, err := http.NewRequest(string(n.Method), n.Url, strings.NewReader(string(reqBody)))
	if err != nil {
		retError := fmt.Errorf("create request failed, error: %s", err.Error())
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			xhlog.Res:     "fail",
			xhlog.Method:  n.Method,
			xhlog.Address: n.Url,
			xhlog.Msg:     retError.Error(),
		})
		return http.StatusBadRequest, nil
	}

	// head处理
	if head != nil {
		for k, v := range *head {
			req.Header[k] = v
		}
	}
	req.Header.Set(xhid.TraceId, ctx.Values().GetString(xhid.TraceId))
	currentSpanId := xhid.MakeSpanId("currentSpanId")
	req.Header.Set(xhid.ParentId, currentSpanId) // 为下游生成一个parent id，并在http请求钟打印该parent id为当前级别的span id
	nextSpanId := xhid.MakeSpanId(n.Url)
	req.Header.Set(xhid.SpanId, nextSpanId) // 为下游新生成span id

	resp, err := httpClient.Do(req)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		})
		return http.StatusBadRequest, nil
	}
	if resp == nil {
		err := fmt.Errorf("response nil")
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		})
		return http.StatusBadRequest, nil
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
			"statusCode":    resp.StatusCode,
		})
		return resp.StatusCode, nil
	}

	xhlog.Info(ctx, xhlog.OPHttpSuccess, map[string]interface{}{
		"spanid":        currentSpanId,
		xhlog.TimeStamp: timestamp,
		xhlog.Res:       "success",
		xhlog.Method:    n.Method,
		xhlog.Address:   n.Url,
		xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
		xhlog.Result:    xhlog.TripSpaceAndReturn(respBody),
		xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		"statusCode":    resp.StatusCode,
	})
	return resp.StatusCode, respBody
}

// SendWithTimeout 发送http请求
func (n *NetHttp) SendWithTimeout(ctx iris.Context, reqBody []byte, head *http.Header, timeoutSec time.Duration) (int, []byte) {
	startTime := time.Now()
	timestamp := xhlog.NowUSec()

	transport := http.Transport{
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * timeoutSec, //超时时间从5s调整到8s，获取设备树接口在组织下设备很多时易超时。
	}
	httpClient := http.Client{
		Timeout:   timeoutSec,
		Transport: &transport,
	}

	// request 准备
	req, err := http.NewRequest(string(n.Method), n.Url, strings.NewReader(string(reqBody)))
	if err != nil {
		retError := fmt.Errorf("create request failed, error: %s", err.Error())
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			xhlog.Res:     "fail",
			xhlog.Method:  n.Method,
			xhlog.Address: n.Url,
			xhlog.Msg:     retError.Error(),
		})
		return http.StatusBadRequest, nil
	}

	// head处理
	if head != nil {
		for k, v := range *head {
			req.Header[k] = v
		}
	}
	req.Header.Set(xhid.TraceId, ctx.Values().GetString(xhid.TraceId))
	currentSpanId := xhid.MakeSpanId("currentSpanId")
	req.Header.Set(xhid.ParentId, currentSpanId) // 为下游生成一个parent id，并在http请求钟打印该parent id为当前级别的span id
	nextSpanId := xhid.MakeSpanId(n.Url)
	req.Header.Set(xhid.SpanId, nextSpanId) // 为下游新生成span id

	resp, err := httpClient.Do(req)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		})
		return http.StatusBadRequest, nil
	}
	if resp == nil {
		err := fmt.Errorf("response nil")
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		})
		return http.StatusBadRequest, nil
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
			"statusCode":    resp.StatusCode,
		})
		return resp.StatusCode, nil
	}

	xhlog.Info(ctx, xhlog.OPHttpSuccess, map[string]interface{}{
		"spanid":        currentSpanId,
		xhlog.TimeStamp: timestamp,
		xhlog.Res:       "success",
		xhlog.Method:    n.Method,
		xhlog.Address:   n.Url,
		xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
		xhlog.Result:    xhlog.TripSpaceAndReturn(respBody),
		xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		"statusCode":    resp.StatusCode,
	})
	return resp.StatusCode, respBody
}

// HttpWithOption 发送http请求
func (n *NetHttp) HttpWithOption(ctx iris.Context, reqBody []byte, head *http.Header, timeoutSec time.Duration) (int, []byte) {
	startTime := time.Now()
	timestamp := xhlog.NowUSec()

	transport := http.Transport{
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * timeoutSec, //超时时间从5s调整到8s，获取设备树接口在组织下设备很多时易超时。
	}
	if n.IsHttps {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // https，则增加额外的判断
	}
	httpClient := http.Client{
		Timeout:   timeoutSec,
		Transport: &transport,
	}

	// request 准备
	req, err := http.NewRequest(string(n.Method), n.Url, strings.NewReader(string(reqBody)))
	if err != nil {
		retError := fmt.Errorf("create request failed, error: %s", err.Error())
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			xhlog.Res:     "fail",
			xhlog.Method:  n.Method,
			xhlog.Address: n.Url,
			xhlog.Msg:     retError.Error(),
		})
		return http.StatusBadRequest, nil
	}

	// head处理
	if head != nil {
		for k, v := range *head {
			req.Header[k] = v
		}
	}
	req.Header.Set(xhid.TraceId, ctx.Values().GetString(xhid.TraceId))
	currentSpanId := xhid.MakeSpanId("currentSpanId")
	req.Header.Set(xhid.ParentId, currentSpanId) // 为下游生成一个parent id，并在http请求钟打印该parent id为当前级别的span id
	nextSpanId := xhid.MakeSpanId(n.Url)
	req.Header.Set(xhid.SpanId, nextSpanId) // 为下游新生成span id

	resp, err := httpClient.Do(req)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		})
		return http.StatusBadRequest, nil
	}
	if resp == nil {
		err := fmt.Errorf("response nil")
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		})
		return http.StatusBadRequest, nil
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPHttpFailure, map[string]interface{}{
			"spanid":        currentSpanId,
			xhlog.TimeStamp: timestamp,
			xhlog.Res:       "fail",
			xhlog.Method:    n.Method,
			xhlog.Address:   n.Url,
			xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
			xhlog.Msg:       err.Error(),
			xhlog.ProcTime:  xhlog.GetProcTime(startTime),
			"statusCode":    resp.StatusCode,
		})
		return resp.StatusCode, nil
	}

	xhlog.Info(ctx, xhlog.OPHttpSuccess, map[string]interface{}{
		"spanid":        currentSpanId,
		xhlog.TimeStamp: timestamp,
		xhlog.Res:       "success",
		xhlog.Method:    n.Method,
		xhlog.Address:   n.Url,
		xhlog.Body:      xhlog.TripSpaceAndReturn(reqBody),
		xhlog.Result:    xhlog.TripSpaceAndReturn(respBody),
		xhlog.ProcTime:  xhlog.GetProcTime(startTime),
		"statusCode":    resp.StatusCode,
	})
	return resp.StatusCode, respBody
}

func NewNetHttp(method Method, url string, isHttps bool) *NetHttp {
	return &NetHttp{
		Method:  method,
		Url:     url,
		IsHttps: isHttps,
	}
}
