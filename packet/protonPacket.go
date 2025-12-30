package packet

import (
	"encoding/json"
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
	payload     []byte // 存储原始payload字节数据，用于后续转换
}

// SetPayload 设置payload数据
func (packet *ProtonPacket[T]) SetPayload(data []byte) {
	packet.payload = data
}

// GetPayload 获取payload数据
func (packet *ProtonPacket[T]) GetPayload() []byte {
	return packet.payload
}

// toJsonObject 将packet中的payload转成传入的结构体类型并返回
func (packet *ProtonPacket[T]) toJsonObject() (T, error) {
	var result T
	if packet.payload == nil {
		return result, nil
	}

	err := json.Unmarshal(packet.payload, &result)
	return result, err
}
