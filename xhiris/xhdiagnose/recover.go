package xhdiagnose

import (
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"runtime"
)

// RecoverFunc 支持defer调用，用于panic的恢复
func RecoverFunc() {
	if err := recover(); err != nil {
		buf := make([]byte, 1024)
		runtime.Stack(buf, true)
		xhlog.AppError("PanicRecover", map[string]interface{}{
			xhlog.ErrorMsg: err,
			"stack":        string(buf),
		})
	}
}
