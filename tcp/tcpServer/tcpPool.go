package tcpServer

import (
	"context"
	"net"
)

type WorkPool struct {
	tasks  chan net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorkPool(size int) *WorkPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &WorkPool{
		tasks:  make(chan net.Conn, size),
		ctx:    ctx,
		cancel: cancel,
	}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return pool
}

func (pool *WorkPool) worker() {
	for {
		select {
		case task, ok := <-pool.tasks:
			if !ok {
				// 通道已关闭
				return
			}
			HandleConnect(task)
		case <-pool.ctx.Done():
			// 上下文被取消，退出worker
			return
		}
	}
}

func (pool *WorkPool) AddTask(task net.Conn) {
	select {
	case pool.tasks <- task:
	case <-pool.ctx.Done():
		// 如果上下文已取消，关闭连接
		err := task.Close()
		if err != nil {
			return
		}
	}
}

func (pool *WorkPool) Close() {
	pool.cancel()
	close(pool.tasks)
}
