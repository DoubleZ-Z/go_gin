package broker

import (
	"go_gin/packetV2"
	"go_gin/tcp/dto"
)

type Broker interface {
	HandlePacket(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) string
	SendAction(stationNo string, action string, actionVar string, content any) string
}

type DefaultBroker struct {
}

func (b *DefaultBroker) HandlePacket(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) string {
	packet
	return ""
}

func (b *DefaultBroker) SendAction(stationNo string, action string, actionVar string, content any) string {
	return ""
}
