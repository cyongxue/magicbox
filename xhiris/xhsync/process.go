package xhsync

import (
	"context"
	"sync"
)

// 处理器
type Process struct {
	Index   ProcessIndex
	EventCh chan GroupEvent
}

func (p *Process) Start(rawCtx context.Context, wg sync.WaitGroup) {
	wg.Add(1)
	// 启动一个独立协程进行处理
	go func(rawCtx context.Context) {
		defer wg.Done()

		for {
			select {
			case ev := <-p.EventCh:
				// 业务处理
				ev.Deal()
			case <-rawCtx.Done():
				// 用于结束协程
				return
			}
		}
	}(rawCtx)

	return
}
