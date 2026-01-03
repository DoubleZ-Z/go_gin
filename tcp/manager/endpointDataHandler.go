package manager

import (
	"fmt"
	"go_gin/packetV2"
	"go_gin/tcp/dto"
	"net/http"
)

func OnRequest(packet packetV2.ProtonPacket[any], connect dto.TcpConnect) packetV2.ProtonPacket[any] {
	header := packet.Header

	var response packetV2.ProtonPacket[any]
	var err error
	handler := GetActionHandler(header.Action, header.ActionVar)
	if handler == nil {
		err = fmt.Errorf("cannot find handler:action:[%s],actionVer:[%s]", header.Action, header.ActionVar)
	} else {
		payload := handler.HandleAction(packet, connect)
		response, err = packetV2.Response(packet, payload, http.StatusOK, "OK")
	}
	if err != nil {
		response, err = packetV2.Response(packet, nil, http.StatusInternalServerError, "handle action["+header.Action+"] err")
	}
	return response
}
