package xhsync

import (
	"context"
	"fmt"
	"sync"
)

type ProcessIndex int

// 运行池
type Engine struct {
	wg     sync.WaitGroup
	rawCtx context.Context
	cancel context.CancelFunc

	poolCap    int
	pool       map[ProcessIndex]*Process
	processCap int
}

var EnginePrt *Engine

// NewGroupEngine 新建一个团购处理引擎
func NewGroupEngine(ctx context.Context, cancel context.CancelFunc, poolCap int, processCap int) error {
	if poolCap < 8 {
		poolCap = 8
	}
	if processCap < 1024 {
		processCap = 1024
	}

	engine := &Engine{
		rawCtx:     ctx,
		cancel:     cancel,
		pool:       make(map[ProcessIndex]*Process, poolCap),
		poolCap:    poolCap,
		processCap: processCap,
	}

	for i := 0; i < poolCap; i++ {
		// 创建process
		process := &Process{
			Index:   ProcessIndex(i),
			EventCh: make(chan GroupEvent, processCap),
		}
		// 启动process
		process.Start(engine.rawCtx, engine.wg)

		engine.pool[process.Index] = process
	}

	EnginePrt = engine
	return nil
}

const (
	DefaultId = 0
)

// GroupEngine.Dispatch 拼团写时间调度处理，保证同一个groupon只会在一个goroutine中进行处理，故同一个groupon的事件处理是串行的
func (g *Engine) Dispatch(id int64, event GroupEvent) error {
	mod := id % int64(g.poolCap)
	index := ProcessIndex(mod)

	if _, had := g.pool[index]; !had {
		return fmt.Errorf("process index=%d no exist", index)
	}

	// 投递给具体process进行处理
	g.pool[index].EventCh <- event

	return nil
}
