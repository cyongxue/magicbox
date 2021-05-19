package xhredis

import (
	"fmt"
	"github.com/cyongxue/magicbox/xhiris/xhid"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12/context"
	"sync"
	"testing"
	"time"
)

func TestCmdNoCtxCtx(t *testing.T) {
	logConf := xhlog.LoggerConf{
		Dir:        "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhredis/logs",
		Prefix:     "test",
		Level:      "info",
		RotateSize: 1 * 1024 * 1024,
		Console:    true,
	}
	if err := xhlog.Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	config := Config{
		Server:      "127.0.0.1:6379",
		Pwd:         "",
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: time.Duration(60) * time.Second,
	}
	engine := NewEngine(&config)

	type DeviceMediaInfo struct {
		PhoneNo   string
		EndTime   string
		BeginTime string
	}
	deviceMediaInfo := DeviceMediaInfo{
		PhoneNo:   "30000276",
		EndTime:   "20500731T084824",
		BeginTime: "20200731T084824",
	}
	res, err := engine.CmdNoCtx("HMSET", "5PPI07001815R_0_dev_5PPI07001815R_20200726T195233", deviceMediaInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	str, err := redis.String(res, err)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(str)
	return
}

func TestRedisEngine_Multi(t *testing.T) {
	logConf := xhlog.LoggerConf{
		Dir:        "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhredis/logs",
		Prefix:     "test",
		Level:      "info",
		RotateSize: 1 * 1024 * 1024,
		Console:    true,
	}
	if err := xhlog.Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	config := Config{
		Server:      "127.0.0.1:6379",
		Pwd:         "",
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: time.Duration(60) * time.Second,
	}
	engine := NewEngine(&config)

	var args []MultiArgs
	args = append(args, MultiArgs{Cmd: "SREM", Key: "snapset", Args: []interface{}{"chengyongxue"}})
	args = append(args, MultiArgs{Cmd: "DEL", Key: "chengyongxue"})
	err := engine.MultiNoCtx(args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	members, err := redis.Strings(engine.cmdNoCtx("SMEMBERS", "snapset"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(members)
	return
}

func TestRedisEngine_Acquire(t *testing.T) {
	logConf := xhlog.LoggerConf{
		Dir:        "/Users/hehui/GolandProjects/src/github.com/cyongxue/magicbox/xhiris/xhredis/logs",
		Prefix:     "test",
		Level:      "info",
		RotateSize: 1 * 1024 * 1024,
		Console:    true,
	}
	if err := xhlog.Init(&logConf); err != nil {
		fmt.Println("dh log init error: " + err.Error())
		return
	}

	config := Config{
		Server:      "127.0.0.1:6379",
		Pwd:         "",
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: time.Duration(60) * time.Second,
	}
	engine := NewEngine(&config)

	irisCtx := context.NewContext(nil)
	irisCtx.Values().Set(xhid.TraceId, xhid.IdDriver(1))
	irisCtx.Values().Set(xhid.SpanId, xhid.MakeSpanId("hello world"))

	var mutex sync.Mutex
	lock, err := engine.Acquire(irisCtx, mutex, "locker", "1", 10)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("lock result=", lock)

	if err := engine.Release(irisCtx, mutex, "locker"); err != nil {
		fmt.Println(err.Error())
		return
	}
}
