package sigserver

import (
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"github.com/gictorbit/peershare/utils"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
)

type PacketBody struct {
	MessageType pb.MessageType
	StatusCode  pb.StatusCode
	Payload     []byte
}

func (pss *PeerShareServer) ReadPacket(conn net.Conn) (*PacketBody, error) {
	buf := make([]byte, utils.PacketMaxByteLength)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	messageType := pb.MessageType(buf[0])
	statusCode := pb.StatusCode(buf[1])
	payload := buf[2:n]
	return &PacketBody{
		MessageType: messageType,
		Payload:     payload,
		StatusCode:  statusCode,
	}, nil
}

func (pss *PeerShareServer) SendResponsePacket(conn net.Conn, packet *PacketBody) error {
	buf := make([]byte, 0)
	buf = append(buf, byte(packet.MessageType))
	buf = append(buf, byte(packet.StatusCode))
	buf = append(buf, packet.Payload...)
	if len(buf) > utils.PacketMaxByteLength {
		return utils.ErrInvalidPacketSize
	}
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func (pss *PeerShareServer) SendResponse(conn net.Conn, msgType pb.MessageType, status pb.StatusCode, msg proto.Message) error {
	respBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	packet := &PacketBody{
		MessageType: msgType,
		Payload:     respBytes,
		StatusCode:  status,
	}
	if e := pss.SendResponsePacket(conn, packet); e != nil {
		return e
	}
	return nil
}
