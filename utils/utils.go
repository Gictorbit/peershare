package utils

import (
	"errors"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"google.golang.org/protobuf/proto"
	"net"
)

const (
	PacketMaxByteLength = 2048
	ServerSocketType    = "tcp"
)

var (
	ErrInvalidPacketSize = errors.New("invalid packet size")
)

type MessageBody[T proto.Message] struct {
	MessageType pb.MessageType
	Payload     []byte
	Message     T
}

func ReadMessageFromConn[T proto.Message](conn net.Conn, message T) (*MessageBody[T], error) {
	buf := make([]byte, PacketMaxByteLength)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	data := buf[:n]
	packet := &MessageBody[T]{
		MessageType: pb.MessageType(data[0]),
		Payload:     data[1:],
	}
	if e := proto.Unmarshal(packet.Payload, message); e != nil {
		return nil, e
	}
	packet.Message = message

	return packet, nil
}
