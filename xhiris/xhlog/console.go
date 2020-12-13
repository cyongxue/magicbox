package xhlog

import (
	"fmt"
	"time"
)

func console(printLevel Level, v ...interface{}) {
	strTime := time.Now().Format(TimeFormat)
	logStr := fmt.Sprintf("%s %s %s", levels[printLevel].RawText, strTime, fmt.Sprint(v...))
	fmt.Println(logStr)
	return
}
