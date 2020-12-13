package xhprotocol

import (
	"github.com/kataras/iris/v12"
)

// BaseResponse 基础返回信息
type BaseResponse struct {
	TraceId string      `json:"trace_id"`
	Code    int         `json:"code"`
	ErrMsg  string      `json:"desc"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	RespCode = "error_code"
	RespData = "response_data"
)

// SetResponse 返回内容
func SetResponse(ctx iris.Context, err ResponseError, data interface{}) {
	ctx.Values().Set(RespCode, err)
	ctx.Values().Set(RespData, data)
	return
}

// SetResponseOK 设置返回OK
func SetResponseOK(ctx iris.Context, data interface{}) {
	ctx.Values().Set(RespCode, OK)
	if data != nil {
		ctx.Values().Set(RespData, data)
	}
	return
}

// SetResponseErr 设置错误
func SetResponseErr(ctx iris.Context, err ResponseError) {
	ctx.Values().Set(RespCode, err)
	return
}
