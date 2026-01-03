package tcpServer

import (
	"fmt"
	"go_gin/tcp/dto"
	"sync"
)

var once sync.Once
var tcpConnectionAbstract *TcpConnectionAbstract

type TcpConnectionAbstractInterface interface {
	AddOpenChannel(conn dto.TcpConnect)
	AddStationChannel(stationNo string, conn dto.TcpConnect)
	RemoveChannel(id string)
}

type TcpConnectionAbstract struct {
	// 缓存所有连接
	ConnectionCache map[string]dto.TcpConnect

	// 缓存已验证连接
	StationConnectionCache map[string]dto.TcpConnect

	// 缓存channel -> station
	ChannelMap map[string]string
}

func GetTcpConnectionAbstract() *TcpConnectionAbstract {
	once.Do(func() {
		tcpConnectionAbstract = &TcpConnectionAbstract{}
	})
	return tcpConnectionAbstract
}

func (t *TcpConnectionAbstract) init() {
	t.ConnectionCache = make(map[string]dto.TcpConnect)
	t.StationConnectionCache = make(map[string]dto.TcpConnect)
	t.ChannelMap = make(map[string]string)
}

func (t *TcpConnectionAbstract) AddOpenChannel(conn dto.TcpConnect) {
	t.ConnectionCache[conn.Id] = conn
}

func (t *TcpConnectionAbstract) AddStationChannel(stationNo string, conn dto.TcpConnect) {
	if connect, ok := t.StationConnectionCache[stationNo]; ok {
		if connect.Id != conn.Id {
			t.RemoveChannel(connect.Id)
		}
	} else {
		t.StationConnectionCache[stationNo] = conn
	}
}

func (t *TcpConnectionAbstract) RemoveChannel(id string) {
	if connect, ok := t.ConnectionCache[id]; ok {
		err := connect.Conn.Close()
		if err != nil {
			fmt.Println(err)
		}
		delete(t.ConnectionCache, id)
		if stationNo, ok := t.ChannelMap[id]; ok {
			delete(t.ChannelMap, id)
			if tcpConnect, ok := t.StationConnectionCache[stationNo]; ok {
				if tcpConnect.Id == id {
					delete(t.StationConnectionCache, stationNo)
				}
			}
		}
	}
}

func (t *TcpConnectionAbstract) GetChannel(stationNo string) (dto.TcpConnect, bool) {
	if connect, ok := t.StationConnectionCache[stationNo]; ok {
		return connect, true
	}
	return dto.TcpConnect{}, false
}
