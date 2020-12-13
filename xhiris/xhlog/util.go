package xhlog

import (
	"strconv"
	"strings"
	"time"
)

// 计算处理时间，开始时间或持续的时间
func GetProcTime(startOrDuration interface{}) string {
	switch val := startOrDuration.(type) {
	case time.Time:
		elapsed := time.Now().Sub(val) / 1e2
		elapsedMs := float64(int64(elapsed)) / 10.0
		return strconv.FormatFloat(elapsedMs, 'f', -1, 64)
	case time.Duration:
		elapsedMs := float64(int64(val/1e2)) / 10.0
		return strconv.FormatFloat(elapsedMs, 'f', -1, 64)
	default:
		return "0.0"
	}
}

// NowMsec 获取的毫秒时间
func NowMsec() int64 {
	return time.Now().UnixNano() / 1e6
}

// NowUSec 获取当前的微秒时间
func NowUSec() int64 {
	return time.Now().UnixNano() / 1e3
}

// TripSpaceAndReturn 为打印body去除回车和空格
func TripSpaceAndReturn(body []byte) string {
	if body == nil {
		return ""
	}
	tmp := strings.Replace(string(body), "  ", "", -1)
	ret := strings.Replace(tmp, "\n", "", -1)
	return ret
}
