package xhlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cyongxue/magicbox/xhiris/xhid"
	"github.com/kataras/iris/v12"
	"reflect"
	"runtime"
	"strings"
)

const (
	RealRemoteIP = "x-real-ip"
)

// 日志格式为：
// [INFO] [2020-06-29T22:26:59.972+0800] [logic/middleware/middlerware_handler.go:64] _request_in||uri=/api/ping?world=hello||from=[::1]:63846||method=GET||proto=HTTP/1.1||traceid=1234567890||spanid=1234567890||args={"world":"hello"}

// buildValueLog 拼接日志全文
// 拼接的内容为： [logic/middleware/middlerware_handler.go:64] _request_in||uri=/api/ping?world=hello||from=[::1]:63846||method=GET||proto=HTTP/1.1||traceid=1234567890||spanid=1234567890||args={"world":"hello"}
func buildValueLog(stackDeep int, dFlag string, ctx iris.Context, v ...interface{}) string {

	_, file, line, ok := runtime.Caller(stackDeep)
	if !ok {
		file = "???"
		line = 0
	} else {
		beginIndex := strings.LastIndex(file, "/") + 1
		if beginIndex <= len(file) {
			file = file[beginIndex:]
		}
	}

	var msgBuffer bytes.Buffer
	msgBuffer.WriteString(fmt.Sprintf("[%s:%d]", file, line))

	// 添加dflag
	{
		msgBuffer.WriteString(" ")
		msgBuffer.WriteString(fmt.Sprintf("op=%s", dFlag))
	}

	// 普通日志，采用当前打印的时间戳
	if dFlag != OPHttpSuccess && dFlag != OPHttpFailure {
		msgBuffer.WriteString(fmt.Sprintf("||%s=%d", TimeStamp, NowUSec()))
	}
	// 上下文件的一些信息，利用Context来完成
	if ctx != nil {
		if dFlag == OPRequestIn {
			msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "uri", ctx.Request().URL.String()))
			if ctx.Request().Header.Get(RealRemoteIP) != "" {
				msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "from", ctx.Request().Header.Get(RealRemoteIP)))
			} else {
				msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "from", ctx.Request().RemoteAddr))
			}
		}
		if dFlag == OPRequestIn || dFlag == OPRequestOut {
			msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "method", ctx.Request().Method))
			msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "proto", ctx.Request().Proto))
		}

		msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "traceid", ctx.Values().GetStringDefault(xhid.TraceId, "")))
		msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "parentid", ctx.Values().GetStringDefault(xhid.ParentId, "")))
		if dFlag != OPHttpFailure && dFlag != OPHttpSuccess {
			msgBuffer.WriteString(fmt.Sprintf("||%s=%s", "spanid", ctx.Values().GetStringDefault(xhid.SpanId, "")))
		}
	}

	// 利用v自定义携带的k-v
	{
		paramSlice := []interface{}{}
		paramSlice = append(paramSlice, v...)
		appendValueToBytesBuffer(&msgBuffer, paramSlice...)
	}

	return msgBuffer.String()
}

// appendValueToBytesBuffer 拼接日志格式
func appendValueToBytesBuffer(stringBuffer *bytes.Buffer, v ...interface{}) {

	for _, item := range v {
		// 如果参数是nil，那么就不输出
		if item == nil {
			continue
		}

		// 是不是空
		itemValue := reflect.ValueOf(item)

		switch itemValue.Kind() {
		case reflect.Map:
			// 是不是空
			if itemValue.IsNil() {
				continue
			}
			// 拿到所有key
			keys := itemValue.MapKeys()
			for _, key := range keys {
				value := itemValue.MapIndex(key)

				for value.Kind() == reflect.Interface ||
					value.Kind() == reflect.Ptr {

					value = value.Elem()
				}

				switch value.Kind() {
				case reflect.Map, reflect.Slice, reflect.Struct:
					jsonString, _ := json.Marshal(value.Interface())
					stringBuffer.WriteString(fmt.Sprintf("||%s=%s", key.Interface(), jsonString))
				case reflect.Invalid:
					stringBuffer.WriteString(fmt.Sprintf("||%s=%v", key.Interface(), nil))
				default:
					stringBuffer.WriteString(fmt.Sprintf("||%s=%s", key.Interface(), fmt.Sprint(value.Interface())))
				}
			}

		case reflect.Struct:
			var logString string
			logBytes, _ := json.Marshal(v)
			logString = fmt.Sprintf("||msg=%s", logBytes)

			stringBuffer.WriteString(logString)

		case reflect.Invalid:
			stringBuffer.WriteString(fmt.Sprintf("||msg=%v", nil))

		default:
			stringBuffer.WriteString(fmt.Sprintf("||msg=%s", fmt.Sprint(item)))
		}
	}

	return
}
