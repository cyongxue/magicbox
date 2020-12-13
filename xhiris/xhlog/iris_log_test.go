package xhlog

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"testing"
	"time"
)

func TestAppWarn(t *testing.T) {
	logConf := LoggerConf{
		Dir:     "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlog/logs",
		Prefix:  "test",
		Level:   "debug",
		Console: true,
	}
	if err := Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	app := iris.New()
	app.Use(recover.New())

	AppWarn(OPUndef, map[string]interface{}{
		"framework": "framework",
		"logger": map[string]string{
			"city":    "hangzhou",
			"country": "china",
		},
	})
	AppDebug(OPRequestIn, map[string]interface{}{
		"framework": "framework",
		"logger": map[string]string{
			"city":    "hangzhou",
			"country": "china",
		},
	})
	AppInfo(OPRequestIn, map[string]interface{}{
		"framework": "framework",
		"logger": map[string]string{
			"city":    "hangzhou",
			"country": "china",
		},
	})
	AppError(OPRequestIn, map[string]interface{}{
		"framework": "framework",
		"logger": map[string]string{
			"city":    "hangzhou",
			"country": "china",
		},
	})
	<-time.After(5 * time.Second)
}

func BenchmarkAppInfo(b *testing.B) {
	logConf := LoggerConf{
		Dir:     "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlog/logs",
		Prefix:  "test",
		Level:   "debug",
		Console: true,
	}
	if err := Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	app := iris.New()
	app.Use(recover.New())

	for i := 0; i < b.N; i++ {
		AppInfo(OPRequestIn, map[string]interface{}{
			"framework": "framework",
			"logger": map[string]string{
				"city":    "hangzhou",
				"country": "china",
			},
			"logConf": logConf,
		})
	}

}

func BenchmarkAppWarn(b *testing.B) {
	logConf := LoggerConf{
		Dir:     "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlog/logs",
		Prefix:  "test",
		Level:   "debug",
		Console: true,
	}

	for i := 0; i < b.N; i++ {
		fmt.Println(logConf)
	}
}
