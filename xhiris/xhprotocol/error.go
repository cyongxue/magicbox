package xhprotocol

// ResponseError 返回的错误信息
type ResponseError struct {
	Code     int
	InnerMsg string
	UserMsg  string
}

// ResponseError.Error 优先返回外部错误
func (e ResponseError) Error() string {
	if e.UserMsg == "" {
		return e.InnerMsg
	}
	return e.UserMsg
}

var (
	OK = ResponseError{Code: 0, InnerMsg: "成功", UserMsg: "success"}

	ErrParams  = ResponseError{Code: 1010, InnerMsg: "参数错误", UserMsg: "params error"}
	ErrUnknown = ResponseError{Code: 9998, InnerMsg: "未知错误", UserMsg: "unknown error"}
	ErrSystem  = ResponseError{Code: 9999, InnerMsg: "系统错误", UserMsg: "system error"}
)
