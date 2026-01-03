package dto

import "net"

type TcpConnect struct {
	Conn net.Conn
	Id   string
}
