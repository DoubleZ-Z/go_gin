package tcpServer

import (
	"fmt"
	"go_gin/models"
	"go_gin/packetV2"
	"go_gin/tcp/broker"
	"go_gin/tcp/dto"
	"net"
	"time"
)

const (
	HEADER = 0x02
	FCS    = 0x03
)

type Manager struct {
	tcpChannelMap map[string]net.TCPAddr
}

func HandleConnect(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("连接关闭错误:", err)
		}
		fmt.Println("连接已关闭并释放资源")
	}()
	connect := dto.TcpConnect{
		Conn: conn,
		Id:   models.GenerateUniqueID(),
	}
	connectionAbstract := GetTcpConnectionAbstract()
	connectionAbstract.AddOpenChannel(connect)
	buffer := make([]byte, 5*1024)
	for {
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			fmt.Println("连接读取超时:", err)
			break
		}

		payload, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("读取数据错误或连接已断开:", err)
			break
		}
		startIndex := -1
		endIndex := -1
		for i, b := range buffer[:payload] {
			if b == HEADER && startIndex == -1 {
				startIndex = i
			} else if b == FCS && startIndex != -1 {
				endIndex = i
				break
			}
		}

		if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
			dataPacket := buffer[startIndex+1 : endIndex]
			fmt.Printf("接收到数据包: %s\n", string(dataPacket))

			// 使用简单数据处理器处理数据包
			defaultBroker := broker.DefaultBroker{}
			packetV2 := packetV2.Deserialize(dataPacket)
			// 使用 DataHandler 处理数据包
			result := defaultBroker.HandlePacket(packetV2, connect)
			fmt.Printf("数据包处理结果: %s\n", result)
			startIndex = -1
			endIndex = -1
		} else if startIndex != -1 && endIndex == -1 {
			fmt.Println("数据包不完整，等待更多数据...")
		}
	}
}
