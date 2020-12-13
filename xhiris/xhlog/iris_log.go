package xhlog

import (
	"github.com/kataras/iris/v12"
)

/*
 * op 预定义
 */
const (
	OPUndef = "_undef"

	// 服务端日志
	OPRequestIn  = "_request_in"
	OPRequestOut = "_request_out"

	// 客户端类的日志flag
	// http客户端日志
	OPHttpSuccess = "_http_success"
	OPHttpFailure = "_http_failure"

	// 存储操作日志
	OPMysql        = "_mysql"
	OPRedisSuccess = "_redis_success"
	OPRedisFailure = "_redis_failure"

	// mq操作日志
	OPMqProduct  = "_mq_product"
	OPMqConsumer = "_mq_consume"

	// todo：其他的OP可以业务支持自定义
)

/*
 * log key
 */
const (
	Args     = "args"
	ProcTime = "proc_time"

	ErrorNo  = "errno"
	ErrorMsg = "errmsg"
	Msg      = "msg"
	Res      = "res"

	Method  = "method"
	Body    = "body"
	Address = "address"
	Result  = "result"

	TimeStamp = "timestamp"
	Kind      = "kind"
)

const (
	stackDeep = 2
)

/*
 * Debug debug级别的日志
 */
func Debug(ctx iris.Context, op string, v ...interface{}) {
	LoggerExp.Debug(buildValueLog(stackDeep, op, ctx, v...))
}

func AppDebug(op string, v ...interface{}) {
	LoggerExp.Debug(buildValueLog(stackDeep, op, nil, v...))
}

/*
 * Info info级别的日志
 */
func Info(ctx iris.Context, op string, v ...interface{}) {
	LoggerExp.Info(buildValueLog(stackDeep, op, ctx, v...))
}

func AppInfo(op string, v ...interface{}) {
	LoggerExp.Info(buildValueLog(stackDeep, op, nil, v...))
}

/*
 * Warn warn级别的日志
 */
func Warn(ctx iris.Context, op string, v ...interface{}) {
	LoggerExp.Warn(buildValueLog(stackDeep, op, ctx, v...))
}

func AppWarn(op string, v ...interface{}) {
	LoggerExp.Warn(buildValueLog(stackDeep, op, nil, v...))
}

/*
 * Error error级别的日志
 */
func Error(ctx iris.Context, op string, v ...interface{}) {
	LoggerExp.Error(buildValueLog(stackDeep, op, ctx, v...))
}

func AppError(op string, v ...interface{}) {
	LoggerExp.Error(buildValueLog(stackDeep, op, nil, v...))
}
