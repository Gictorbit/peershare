package client

import (
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"github.com/gictorbit/peershare/utils"
	"google.golang.org/protobuf/proto"
	"net"
)

type PacketBody struct {
	MessageType pb.MessageType
	StatusCode  pb.StatusCode
	Payload     []byte
}

func (pc *PeerClient) ReadPacket(conn net.Conn) (*PacketBody, error) {
	buf := make([]byte, utils.PacketMaxByteLength)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return &PacketBody{
		MessageType: pb.MessageType(buf[0]),
		StatusCode:  pb.StatusCode(buf[1]),
		Payload:     buf[2:n],
	}, nil
}

func (pc *PeerClient) SendRequestPacket(packet *PacketBody) error {
	buf := make([]byte, 0)
	buf = append(buf, byte(packet.MessageType))
	buf = append(buf, packet.Payload...)
	if len(buf) > utils.PacketMaxByteLength {
		return utils.ErrInvalidPacketSize
	}
	if _, err := pc.conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func (pc *PeerClient) SendRequest(msgType pb.MessageType, msg proto.Message) error {
	respBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	packet := &PacketBody{
		MessageType: msgType,
		Payload:     respBytes,
	}
	if e := pc.SendRequestPacket(packet); e != nil {
		return e
	}
	return nil
}
