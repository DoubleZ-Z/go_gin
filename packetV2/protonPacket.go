package packetV2

import (
	"encoding/json"
	"go_gin/util"
	"time"
)

const (
	TYPE_REQUEST   = "request"
	TYPE_RESPONSE  = "response"
	REASON_COMMAND = "command"
	PROTOCOL_VER_2 = 2
)

type ProtonPacketHeader struct {
	Action    string
	ActionVar string
	Trace     string
	Priority  string
	Timestamp string
	Sign      string
}

type ProtonPacketExt struct {
	AppType     string
	AppVer      string
	Station     string
	Timestamp   string
	ResCode     int
	ResMsg      string
	WorksheetNo string
	AccountNo   string
}

type ProtonPacket[T any] struct {
	Type        string
	Reason      string
	ProtocolVar int
	Header      ProtonPacketHeader
	Ext         ProtonPacketExt
	payload     any
}

func CheckSign(packet ProtonPacket[any]) bool {
	sign := Sign(packet, "levent8421")
	if sign == packet.Header.Sign {
		return true
	}
	return false
}

func Deserialize(payloadV2 []byte) ProtonPacket[any] {

	var packet ProtonPacket[any]
	// 首先尝试解析整个字符串为JSON到ProtonPacket结构
	// 尝试解析为ProtonPacket结构
	// 由于我们不知道具体的payload类型，使用any作为泛型参数
	var tempPacket struct {
		Type        string             `json:"type"`
		Reason      string             `json:"reason"`
		ProtocolVar int                `json:"protocolVar"`
		Header      ProtonPacketHeader `json:"header"`
		Ext         ProtonPacketExt    `json:"ext"`
		Payload     json.RawMessage    `json:"payload"`
	}

	err := json.Unmarshal(payloadV2, &tempPacket)
	if err != nil {
		// TODO 解析失败
		return ProtonPacket[any]{
			Type:        "error",
			Reason:      "json unmarshal error",
			ProtocolVar: 0,
			Header:      ProtonPacketHeader{},
			Ext:         ProtonPacketExt{},
			payload:     nil,
		}
	}

	//
	packet.Type = tempPacket.Type
	packet.Reason = tempPacket.Reason
	packet.ProtocolVar = tempPacket.ProtocolVar
	packet.Header = tempPacket.Header
	packet.Ext = tempPacket.Ext
	packet.SetPayload(tempPacket.Payload)

	return packet
}

// SetPayload 设置payload数据
func (packet *ProtonPacket[T]) SetPayload(data any) {
	packet.payload = data
}

// GetPayload 获取payload数据
func (packet *ProtonPacket[T]) GetPayload() any {
	return packet.payload
}

// toJsonObject 将packet中的payload转成传入的结构体类型并返回
func (packet *ProtonPacket[T]) toJsonObject() (T, error) {
	var result T
	if packet.payload == nil {
		return result, nil
	}
	bytes, err := json.Marshal(packet.payload)
	if err != nil {
		return result, err
	}
	err2 := json.Unmarshal(bytes, &result)
	if err2 != nil {
		return result, err2
	}
	return result, err
}

func Response(request ProtonPacket[any], payload any, statusCode int, message string) (ProtonPacket[any], error) {
	response := ProtonPacket[any]{}
	response.SetPayload(payload)
	response.Type = TYPE_RESPONSE
	response.ProtocolVar = PROTOCOL_VER_2

	ext := ProtonPacketExt{
		ResCode: statusCode,
		ResMsg:  message,
	}
	response.Ext = ext
	header := ProtonPacketHeader{
		Action:    request.Header.Action,
		ActionVar: request.Header.ActionVar,
		Trace:     request.Header.Trace,
		Priority:  request.Header.Priority,
		Timestamp: util.UnixToTimestampString(time.Now().Unix()),
	}
	response.Header = header
	return response, nil
}
