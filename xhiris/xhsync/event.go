package xhsync

const (
	ResponseChannelCap = 10
)

type GroupEvent interface {
	Deal()
	SendResponseErr(err error)
}
