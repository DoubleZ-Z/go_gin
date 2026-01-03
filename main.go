package main

import (
	"context"
	"errors"
	"go_gin/interceptor"
	"go_gin/routers"
	"go_gin/tcp/tcpServer"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// ServiceManager 管理多个服务的生命周期
type ServiceManager struct {
	services []Service
	wg       sync.WaitGroup
}

// Service 定义服务接口
type Service interface {
	Start() error
	Stop(ctx context.Context) error
	Name() string
}

// HTTPService HTTP服务实现
type HTTPService struct {
	addr   string
	engine *gin.Engine
	server *http.Server
}

func NewHTTPService(addr string) *HTTPService {
	engine := gin.Default()
	engine.Use(interceptor.Log, routers.TimeCost)
	routers.OpenRouters(engine)
	routers.TokenRouters(engine)

	// 创建标准的http.Server
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	return &HTTPService{
		addr:   addr,
		engine: engine,
		server: server,
	}
}

func (h *HTTPService) Start() error {
	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			println("HTTP服务启动/运行失败:", err.Error())
		}
	}()
	return nil
}

func (h *HTTPService) Stop(ctx context.Context) error {
	if h.server != nil {
		return h.server.Shutdown(ctx)
	}
	return nil
}

func (h *HTTPService) Name() string {
	return "HTTP Service"
}

// TCPService TCP服务实现
type TCPService struct {
	addr     string
	listener net.Listener
	workPool *tcpServer.WorkPool
	stopChan chan struct{}
}

func NewTCPService(addr string) *TCPService {
	return &TCPService{
		addr:     addr,
		stopChan: make(chan struct{}, 1),
	}
}

func (t *TCPService) Start() error {
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	t.listener = listener

	t.workPool = tcpServer.NewWorkPool(10)

	go func() {
		for {
			select {
			case <-t.stopChan:
				return
			default:
				// 非阻塞检查停止信号后，继续接受连接
			}

			connection, err := t.listener.Accept()
			if err != nil {
				// 检查是否是关闭信号
				select {
				case <-t.stopChan:
					return
				default:
					println("接受连接错误:", err.Error())
					continue
				}
			}
			t.workPool.AddTask(connection)
		}
	}()

	return nil
}

func (t *TCPService) Stop(ctx context.Context) error {
	close(t.stopChan)
	if t.listener != nil {
		err := t.listener.Close()
		if err != nil {
			return err
		}
	}
	if t.workPool != nil {
		t.workPool.Close()
	}
	return nil
}

func (t *TCPService) Name() string {
	return "TCP Service"
}

func (sm *ServiceManager) AddService(service Service) {
	sm.services = append(sm.services, service)
}

func (sm *ServiceManager) StartAll() error {
	for _, service := range sm.services {
		if err := service.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (sm *ServiceManager) StopAll() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, service := range sm.services {
		if err := service.Stop(ctx); err != nil {
			// 记录错误但继续关闭其他服务
			println("停止服务错误", service.Name(), ":", err.Error())
		} else {
			println("服务已停止", service.Name())
		}
	}
}

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	// 创建服务管理器
	manager := &ServiceManager{}

	httpService := NewHTTPService(":8080") // 默认端口，也可以通过环境变量配置
	manager.AddService(httpService)

	tcpService := NewTCPService(":11402")
	manager.AddService(tcpService)

	// 启动所有服务
	if err := manager.StartAll(); err != nil {
		panic(err)
	}

	println("所有服务已启动")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	println("正在关闭服务...")
	manager.StopAll()
}
