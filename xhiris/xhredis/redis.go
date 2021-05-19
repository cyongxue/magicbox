package xhredis

import (
	"encoding/json"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"sync"
	"time"
)

// https://www.cnblogs.com/ricklz/p/9562722.html

func (r *RedisEngine) Get(ctx iris.Context, key string) (string, error) {
	resp, err := redis.String(r.CmdWithCtx(ctx, "GET", key))
	if err != nil {
		return "", err
	}
	return resp, nil
}

type MultiArgs struct {
	Cmd  string
	Key  string
	Args []interface{}
}

// RedisEngine.Multi 事务处理
func (r *RedisEngine) Multi(ctx iris.Context, multiArgs []MultiArgs) error {
	startTime := time.Now()
	client := r.pool.Get()
	defer client.Close()

	client.Send("MULTI")
	for _, item := range multiArgs {
		redisArgs := r.redisArgs(item.Key, item.Args...)
		client.Send(item.Cmd, redisArgs...)
	}
	_, err := client.Do("EXEC")
	if err != nil {
		xhlog.Error(ctx, xhlog.OPRedisFailure, map[string]interface{}{
			xhlog.Args:     multiArgs,
			xhlog.ErrorMsg: err.Error(),
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return err
	}

	xhlog.Info(ctx, xhlog.OPRedisSuccess, map[string]interface{}{
		xhlog.Args:     multiArgs,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return nil
}

func (r *RedisEngine) MultiNoCtx(multiArgs []MultiArgs) error {
	startTime := time.Now()
	client := r.pool.Get()
	defer client.Close()

	client.Send("MULTI")
	for _, item := range multiArgs {
		redisArgs := r.redisArgs(item.Key, item.Args...)
		client.Send(item.Cmd, redisArgs...)
	}
	_, err := client.Do("EXEC")
	dataArgs, _ := json.Marshal(multiArgs)
	if err != nil {
		xhlog.AppError(xhlog.OPRedisFailure, map[string]interface{}{
			xhlog.Args:     string(dataArgs),
			xhlog.ErrorMsg: err.Error(),
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return err
	}

	xhlog.AppInfo(xhlog.OPRedisSuccess, map[string]interface{}{
		xhlog.Args:     string(dataArgs),
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return nil
}

func (r *RedisEngine) Acquire(ctx iris.Context, mutex sync.Mutex, key string, value string, second int) (bool, error) {
	startTime := time.Now()

	mutex.Lock()
	defer mutex.Unlock()

	redisArgs := r.redisArgs(key, value, "EX", second, "NX")
	res, err := r.cmdWithCtx(ctx, "SET", redisArgs...)
	result, err := redis.String(res, err)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPRedisFailure, map[string]interface{}{
			xhlog.Res:      "fail",
			"cmd":          "SET",
			xhlog.Args:     redisArgs,
			xhlog.ErrorMsg: err.Error(),
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return false, err
	}

	xhlog.Info(ctx, xhlog.OPRedisFailure, map[string]interface{}{
		xhlog.Res:      "success",
		"cmd":          "SET",
		xhlog.Args:     redisArgs,
		xhlog.Result:   result,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	if result == "" {
		return false, nil
	}
	return true, nil
}

func (r *RedisEngine) Release(ctx iris.Context, mutex sync.Mutex, key string) error {
	startTime := time.Now()

	mutex.Lock()
	defer mutex.Unlock()

	redisArgs := r.redisArgs(key)
	_, err := r.cmdWithCtx(ctx, "DEL", redisArgs...)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPRedisFailure, map[string]interface{}{
			xhlog.Res:      "fail",
			"cmd":          "DEL",
			xhlog.Args:     redisArgs,
			xhlog.ErrorMsg: err.Error(),
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return err
	}

	xhlog.Info(ctx, xhlog.OPRedisFailure, map[string]interface{}{
		xhlog.Res:      "success",
		"cmd":          "DEL",
		xhlog.Args:     redisArgs,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return nil
}
