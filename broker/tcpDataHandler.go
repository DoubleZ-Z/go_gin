package broker

type Broker interface {
	HandlePacket(packet []byte, connId string) string
	SendAction(stationNo string, action string, actionVar string, content any) string
}

type DefaultBroker struct {
}

func (b *DefaultBroker) HandlePacket(packet []byte, connId string) string {
	return ""
}

func (b *DefaultBroker) SendAction(stationNo string, action string, actionVar string, content any) string {
	return ""
}
