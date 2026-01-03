package handler

import (
	"fmt"
	"go_gin/packetV2"
	"go_gin/tcp/dto"
	"go_gin/util"
	"time"
)

type HeartbeatHandler struct {
}

func (h *HeartbeatHandler) HandleAction(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) any {
	fmt.Printf("HeartbeatHandler ======================= [%s]", util.UnixToTimestampString(time.Now().Unix()))
	return nil
}
