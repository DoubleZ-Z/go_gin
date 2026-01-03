package tcpServer

import (
	"fmt"
	"go_gin/packetV2"
	"go_gin/tcp/broker"
	"go_gin/tcp/dto"
	"go_gin/util"
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
		Id:   util.GenerateUniqueID(),
	}

	// 确保连接抽象已初始化
	connectionAbstract := GetTcpConnectionAbstract()
	connectionAbstract.AddOpenChannel(connect)

	// 缓冲区和累积数据
	buffer := make([]byte, 5*1024)
	var dataBuffer []byte // 用于累积不完整的数据包

	for {
		// 设置读取超时
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			fmt.Println("设置连接读取超时失败:", err)
			break
		}

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("读取数据错误或连接已断开:", err)
			break
		}

		// 将新读取的数据追加到累积缓冲区
		dataBuffer = append(dataBuffer, buffer[:n]...)

		// 解析数据包
		processed := 0
		for {
			startIndex := -1
			endIndex := -1

			// 寻找数据包的开始和结束位置
			for i := processed; i < len(dataBuffer); i++ {
				if dataBuffer[i] == HEADER && startIndex == -1 {
					startIndex = i
				} else if dataBuffer[i] == FCS && startIndex != -1 {
					endIndex = i
					break
				}
			}

			// 如果找到了完整的数据包
			if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
				dataPacket := dataBuffer[startIndex+1 : endIndex]
				fmt.Printf("接收到数据包: %s\n", string(dataPacket))

				// 使用数据处理器处理数据包
				defaultBroker := broker.DefaultBroker{}
				packet := packetV2.Deserialize(dataPacket)

				// 使用 DataHandler 处理数据包
				result, err := defaultBroker.HandlePacket(packet, connect)
				if err != nil {
					fmt.Printf("数据包处理错误: %v\n", err)
				} else {
					fmt.Printf("数据包处理结果: %s\n", result)
				}

				// 移动已处理的数据位置
				processed = endIndex + 1

				// 重置查找索引
				startIndex = -1
				endIndex = -1
			} else {
				// 没有找到完整的数据包，将已处理的数据部分移除
				if processed > 0 {
					dataBuffer = dataBuffer[processed:]
					processed = 0
				}
				break // 等待更多数据
			}
		}
	}

	// 连接断开时从连接管理器中移除
	connectionAbstract.RemoveChannel(connect.Id)
}
