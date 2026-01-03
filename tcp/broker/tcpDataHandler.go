package broker

import (
	"fmt"
	"go_gin/packetV2"
	"go_gin/tcp/dto"
	"go_gin/tcp/manager"
)

type Broker interface {
	HandlePacket(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) string
	SendAction(stationNo string, action string, actionVar string, content any) string
}

type DefaultBroker struct {
}

func (b *DefaultBroker) HandlePacket(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) (string, error) {
	ok := packetV2.CheckSign(packet)
	if !ok {
		return "", fmt.Errorf("[%d] ,签名错误", packet.Ext.Station)
	}
	switch packet.Type {
	case packetV2.TYPE_REQUEST:
		manager.OnRequest(packet, connect)
	}
	return "", nil
}

func (b *DefaultBroker) SendAction(stationNo string, action string, actionVar string, content any) string {
	return ""
}
