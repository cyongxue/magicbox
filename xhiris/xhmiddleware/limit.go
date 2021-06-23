package dhmiddleware

import (
	"github.com/cyongxue/magicbox/xhiris/xhdiagnose"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	"github.com/kataras/iris/v12"
	"sync"
	"time"
)

// todo: 其他同样类似算法实现
// https://mp.weixin.qq.com/s/5wPpHi8wwaGjen71qXon_A

type LimitUtil struct {
	limitNumber int64 // 触发受限阈值
	minSafeTime int64 // 最小安全期
	limitTime   int64 // 受限时长

	cacheMutex sync.RWMutex
	cacheMap   map[string][]int64 // key ----- [access count, first access time]
	limitMutex sync.RWMutex
	limitMap   map[string]int64 // key ----- release time
}

var LimitEngine *LimitUtil

func InitLimitUtil(limitNumber int, minSafeTime int, limitTime int) {
	LimitEngine = &LimitUtil{
		limitNumber: int64(limitNumber),
		minSafeTime: int64(minSafeTime),
		limitTime:   int64(limitTime),
		cacheMutex:  sync.RWMutex{},
		cacheMap:    make(map[string][]int64),
		limitMutex:  sync.RWMutex{},
		limitMap:    make(map[string]int64),
	}

	go LimitEngine.intervalClearMap()
}

func (l *LimitUtil) intervalClearMap() {
	defer xhdiagnose.RecoverFunc()

	ticker := time.NewTicker(time.Duration(l.minSafeTime) * time.Second)
	for {
		select {
		case <-ticker.C:
			l.cacheMutex.Lock()
			l.cacheMap = make(map[string][]int64)
			l.cacheMutex.Unlock()
		}
	}
}

func (l *LimitUtil) filterLimitedMap() {
	l.limitMutex.Lock()
	defer l.limitMutex.Unlock()

	nowSec := time.Now().Unix()
	for k, v := range l.limitMap {
		if v <= nowSec {
			delete(l.limitMap, k)
		}
	}
}

func (l *LimitUtil) isLimit(ctx iris.Context, key string) bool {
	l.limitMutex.RLock()
	defer l.limitMutex.RUnlock()

	if _, has := l.limitMap[key]; has {
		return true
	}
	return false
}

func (l *LimitUtil) initKeyNumber(ctx iris.Context, key string) {
	info := make([]int64, 2)
	info[0] = 0
	info[1] = time.Now().Unix()
	l.cacheMap[key] = info
}

func (l *LimitUtil) CheckKeyLimit(ctx iris.Context, key string) bool {
	if len(key) == 0 {
		return false
	}

	l.filterLimitedMap()
	if l.isLimit(ctx, key) {
		xhlog.Warn(ctx, "CheckKeyLimit", map[string]interface{}{
			"key": key,
			"msg": "is limited",
		})
		return true
	}

	l.cacheMutex.Lock()
	defer l.cacheMutex.Unlock()
	if _, has := l.cacheMap[key]; has {
		info := l.cacheMap[key]
		info[0] = info[0] + 1
		if info[0] > l.limitNumber {
			firstAccessTime := info[1]
			now := time.Now().Unix()
			if now-firstAccessTime <= l.minSafeTime {
				l.limitMutex.Lock()
				l.limitMap[key] = now + l.limitTime
				l.limitMutex.Unlock()
				xhlog.Warn(ctx, "CheckKeyLimit", map[string]interface{}{
					"key":             key,
					"now":             now,
					"firstAccessTime": firstAccessTime,
					"duration":        now - firstAccessTime,
					"minSafeTime":     l.minSafeTime,
					"accessTime":      info[0],
					"msg":             "is limited",
				})
			} else {
				l.initKeyNumber(ctx, key)
				xhlog.Info(ctx, "CheckKeyLimit", map[string]interface{}{
					"key":             key,
					"firstAccessTime": time.Now().Unix(),
					"minSafeTime":     l.minSafeTime,
					"msg":             "reset count",
				})
			}
		}
	} else {
		l.initKeyNumber(ctx, key)
		xhlog.Info(ctx, "CheckKeyLimit", map[string]interface{}{
			"key":             key,
			"firstAccessTime": time.Now().Unix(),
		})
	}

	return false
}
