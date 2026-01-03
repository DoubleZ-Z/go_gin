package manager

import (
	"go_gin/packetV2"
	"go_gin/tcp/dto"
	"go_gin/tcp/handler"
)

type ActionEndpointInterface interface {
	HandleAction(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) any
}

type ActionEndpoint struct {
	Action        string
	ActionVersion string
}

const (
	HEARTBEAT = "heartbeat.upload"
)

const (
	ACTION_VERSION_1 = "V01"
	ACTION_VERSION_2 = "V02"
)

func init() {
	RegisterAction(HEARTBEAT, ACTION_VERSION_1, &handler.HeartbeatHandler{})
}

var handlerMap = make(map[string]map[string]ActionEndpointInterface)

func RegisterAction(action string, actionVersion string, handler ActionEndpointInterface) {
	if _, ok := handlerMap[action]; !ok {
		handlerMap[action] = make(map[string]ActionEndpointInterface)
	}
	handlerMap[action][actionVersion] = handler
}

func GetActionHandler(action string, actionVersion string) ActionEndpointInterface {
	if actionHandlers, ok := handlerMap[action]; ok {
		if handler, ok := actionHandlers[actionVersion]; ok {
			return handler
		}
	}
	return nil
}
