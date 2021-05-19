package dhmiddleware

import (
	"github.com/cyongxue/magicbox/xhiris/xhid"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/cyongxue/magicbox/xhiris/xhprotocol"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// JSONResponseDone json方式返回的Done封装
func JSONResponseDone() context.Handler {
	return func(ctx iris.Context) {
		var response = xhprotocol.BaseResponse{
			TraceId: ctx.Values().GetString(xhid.TraceId),
		}

		err, ok := ctx.Values().Get(xhprotocol.RespCode).(error)
		if ok {
			if responseErr, ok := err.(xhprotocol.ResponseError); ok {
				response.Code = responseErr.Code
				response.ErrMsg = responseErr.Error()
				// 不论结果如何，都尝试返回data
				if data := ctx.Values().Get(xhprotocol.RespData); data != nil {
					response.Data = data
				}
			} else {
				// 非自定义错误统一返回系统错误
				response.Code = xhprotocol.ErrUnknown.Code
				response.ErrMsg = err.Error()
			}
		} else {
			response.Code = xhprotocol.ErrSystem.Code
			response.ErrMsg = xhprotocol.ErrSystem.Error()
		}

		// 打印request out的日志
		xhlog.Info(ctx, xhlog.OPRequestOut, map[string]interface{}{
			xhlog.Result:   response,
			xhlog.ErrorNo:  response.Code,
			xhlog.ProcTime: xhlog.GetProcTime(ctx.Values().Get(RequestStartTime)),
		})
		ctx.JSON(response)
		return
	}
}

// StatusFoundDone 302 重定向的返回
func StatusFoundDone() context.Handler {
	return func(ctx iris.Context) {
		xhlog.Info(ctx, xhlog.OPRequestOut, map[string]interface{}{
			xhlog.ProcTime: xhlog.GetProcTime(ctx.Values().Get(RequestStartTime)),
			"302":          ctx.ResponseWriter().Header(),
		})
		ctx.StatusCode(iris.StatusFound)
	}
}

// TextResponseDone 如果是跨域请求，直接返回支持跨域
func TextResponseDone() context.Handler {
	return func(ctx iris.Context) {
		data := ctx.Values().GetDefault(xhprotocol.RespData, nil)
		if data != nil {
			xhlog.Info(ctx, xhlog.OPRequestOut, map[string]interface{}{
				xhlog.Result:   data,
				xhlog.ProcTime: xhlog.GetProcTime(ctx.Values().Get(RequestStartTime)),
			})
		} else {
			xhlog.Info(ctx, xhlog.OPRequestOut, map[string]interface{}{
				xhlog.ProcTime: xhlog.GetProcTime(ctx.Values().Get(RequestStartTime)),
			})
		}
		if _, ok := data.(string); ok {
			ctx.WriteString(data.(string))
		}
		ctx.StatusCode(iris.StatusOK)
		return
	}
}
