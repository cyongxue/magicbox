package xhlog

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	logConf := LoggerConf{
		Dir:        "Z:\\Goland\\src\\gitee.com\\yongxue\\magicbox\\main\\logs",
		Prefix:     "test",
		Level:      "info",
		RotateSize: 1 * 1024 * 1024,
	}
	if err := Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	//dhlog.LoggerExp.Error("error error error error error error error error error error error error error")
	go func() {
		for {
			LoggerExp.Error("error error error error error error error error error error error error error")
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func() {
		for {
			LoggerExp.Warn("warn warn warn warn warn warn warn warn warn warn warn warn warn")
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func() {
		for {
			LoggerExp.Info("info info info info info info info info info info info info info")
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func() {
		for {
			LoggerExp.Debug("debug debug debug debug debug debug debug debug debug debug debug debug debug")
			time.Sleep(10 * time.Millisecond)
		}
	}()
	time.Sleep(30 * time.Second)
}

func BenchmarkLogger_Info(b *testing.B) {
	logConf := LoggerConf{
		Dir:        "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhlog/logs",
		Prefix:     "test",
		Level:      "info",
		RotateSize: 1 * 1024 * 1024,
	}
	if err := Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		LoggerExp.Info("info info info info info info info info info info info info infoinfo info info info info info info info info info info info info")
	}
}
