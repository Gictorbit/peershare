package utils

import (
	"encoding/json"
	"errors"
	api "github.com/gictorbit/peershare/api"
	"net"
)

const (
	PacketMaxByteLength = 2048
	ServerSocketType    = "tcp"
)

var (
	ErrInvalidPacketSize = errors.New("invalid packet size")
)

type MessageBody[T any] struct {
	MessageType api.MessageType
	Payload     []byte
	Message     T
}

func ReadMessageFromConn[T any](conn net.Conn, message T) (*MessageBody[T], error) {
	buf := make([]byte, PacketMaxByteLength)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	data := buf[:n]
	packet := &MessageBody[T]{
		MessageType: api.MessageType(data[0]),
		Payload:     data[1:],
	}
	if e := json.Unmarshal(packet.Payload, message); e != nil {
		return nil, e
	}
	packet.Message = message
	return packet, nil
}
