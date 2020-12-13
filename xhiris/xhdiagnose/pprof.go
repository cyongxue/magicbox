package xhdiagnose

import (
	"net/http"
	_ "net/http/pprof"
	"strings"
)

// StartPPROF 启动pprof诊断功能
// addr，一般采用127.0.0.1:port
func StartPPROF(enable bool, addr string) {
	if enable {
		if strings.Contains(addr, ":") {
			go func() {
				_ = http.ListenAndServe(addr, nil)
			}()
		} else {
			go func() {
				_ = http.ListenAndServe("127.0.0.1:"+addr, nil)
			}()
		}
	}
}
