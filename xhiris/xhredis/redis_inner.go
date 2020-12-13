package xhredis

import (
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"strings"
	"time"
)

type RedisEngine struct {
	pool *redis.Pool
}

type Config struct {
	Server      string
	Pwd         string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

// NewEngine 初始化打开redis
func NewEngine(config *Config) *RedisEngine {
	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: config.IdleTimeout,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", config.Server)
			if err != nil {
				xhlog.AppError("RedisInit", map[string]interface{}{
					xhlog.Args: config.Server,
					xhlog.Msg:  err.Error(),
				})
			} else {
				if config.Pwd != "" {
					_, err = conn.Do("AUTH", config.Pwd)
					if err != nil {
						_ = conn.Close()
						conn = nil
						xhlog.AppError("RedisInit", map[string]interface{}{
							xhlog.Args: config.Server,
							xhlog.Msg:  err.Error(),
						})
					}
				}
			}
			return
		},
	}

	redisEngine := &RedisEngine{
		pool: pool,
	}
	return redisEngine
}

func (r *RedisEngine) cmdNoCtx(cmd string, args ...interface{}) (interface{}, error) {
	client := r.pool.Get()
	result, err := client.Do(cmd, args...)
	if err != nil {
		client.Close()
		return nil, err
	}

	client.Close()
	return result, nil
}

func (r *RedisEngine) cmdWithCtx(ctx iris.Context, cmd string, args ...interface{}) (interface{}, error) {
	client := r.pool.Get()
	result, err := client.Do(cmd, args...)
	if err != nil {
		client.Close()
		return nil, err
	}

	client.Close()
	return result, nil
}

func (r *RedisEngine) redisArgs(key string, args ...interface{}) redis.Args {
	redisArgs := redis.Args{}
	redisArgs = redisArgs.Add(key)
	if args != nil && len(args) > 0 {
		for _, a := range args {
			redisArgs = redisArgs.AddFlat(a)
		}
	}
	return redisArgs
}

// CmdNoCtx 不带context的操作
func (r *RedisEngine) CmdNoCtx(cmd string, key string, args ...interface{}) (interface{}, error) {
	startTime := time.Now()

	redisArgs := r.redisArgs(key, args...)
	result, err := r.cmdNoCtx(cmd, redisArgs...)
	if err != nil {
		time.Sleep(10 * time.Millisecond)
		result, e := r.cmdNoCtx(cmd, redisArgs...)
		if e != nil {
			xhlog.AppWarn(xhlog.OPRedisFailure, map[string]interface{}{
				"cmd":          cmd,
				"key":          key,
				xhlog.Args:     args,
				xhlog.ErrorMsg: err.Error(),
				xhlog.ProcTime: xhlog.GetProcTime(startTime),
			})
			return nil, e
		}

		if strings.ToUpper(cmd) == "KEYS" {
			xhlog.AppInfo(xhlog.OPRedisSuccess, map[string]interface{}{
				"cmd":          cmd,
				"key":          key,
				xhlog.Args:     args, // 避免出现打印超大日志，做了keys命令的取舍
				xhlog.ProcTime: xhlog.GetProcTime(startTime),
			})
		} else {
			xhlog.AppInfo(xhlog.OPRedisSuccess, map[string]interface{}{
				"cmd":          cmd,
				"key":          key,
				xhlog.Args:     args,
				xhlog.Result:   result,
				xhlog.ProcTime: xhlog.GetProcTime(startTime),
			})
		}
		return result, nil
	}

	if strings.ToUpper(cmd) == "KEYS" {
		xhlog.AppInfo(xhlog.OPRedisSuccess, map[string]interface{}{
			"cmd":          cmd,
			"key":          key,
			xhlog.Args:     args, // 避免出现打印超大日志，做了keys命令的取舍
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
	} else {
		xhlog.AppInfo(xhlog.OPRedisSuccess, map[string]interface{}{
			"cmd":          cmd,
			"key":          key,
			xhlog.Args:     args,
			xhlog.Result:   result,
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
	}
	return result, nil
}

// CmdWithCtx 带context的操作
func (r *RedisEngine) CmdWithCtx(ctx iris.Context, cmd string, key string, args ...interface{}) (interface{}, error) {
	startTime := time.Now()

	redisArgs := r.redisArgs(key, args...)
	result, err := r.cmdWithCtx(ctx, cmd, redisArgs...)
	if err != nil {
		time.Sleep(10 * time.Millisecond)
		result, e := r.cmdWithCtx(ctx, cmd, redisArgs...)
		if e != nil {
			xhlog.Warn(ctx, xhlog.OPRedisFailure, map[string]interface{}{
				"cmd":          cmd,
				"key":          key,
				xhlog.Args:     args,
				xhlog.ErrorMsg: err.Error(),
				xhlog.ProcTime: xhlog.GetProcTime(startTime),
			})
			return nil, e
		}

		if strings.ToUpper(cmd) == "KEYS" {
			xhlog.Info(ctx, xhlog.OPRedisSuccess, map[string]interface{}{
				"cmd":          cmd,
				"key":          key,
				xhlog.Args:     args, // 避免出现打印超大日志，做了keys命令的取舍
				xhlog.ProcTime: xhlog.GetProcTime(startTime),
			})
		} else {
			xhlog.Info(ctx, xhlog.OPRedisSuccess, map[string]interface{}{
				"cmd":          cmd,
				"key":          key,
				xhlog.Args:     args,
				xhlog.Result:   result,
				xhlog.ProcTime: xhlog.GetProcTime(startTime),
			})
		}
		return result, nil
	}

	if strings.ToUpper(cmd) == "KEYS" {
		xhlog.Info(ctx, xhlog.OPRedisSuccess, map[string]interface{}{
			"cmd":          cmd,
			"key":          key,
			xhlog.Args:     args, // 避免出现打印超大日志，做了keys命令的取舍
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
	} else {
		xhlog.Info(ctx, xhlog.OPRedisSuccess, map[string]interface{}{
			"cmd":          cmd,
			"key":          key,
			xhlog.Args:     args,
			xhlog.Result:   result,
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
	}
	return result, nil
}
