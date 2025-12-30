package tcpServer

import "net"

type WorkPool struct {
	tasks chan net.Conn
}

func NewWorkPool(size int) *WorkPool {
	pool := &WorkPool{
		tasks: make(chan net.Conn, size),
	}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return pool
}

func (pool *WorkPool) worker() {
	for task := range pool.tasks {
		HandleConnect(task)
	}
}

func (pool *WorkPool) AddTask(task net.Conn) {
	pool.tasks <- task
}
