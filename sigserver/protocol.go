package sigserver

import (
	"encoding/json"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"io"
	"net"
)

type PacketBody struct {
	MessageType api.MessageType
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
	messageType := api.MessageType(buf[0])
	payload := buf[1:n]
	return &PacketBody{
		MessageType: messageType,
		Payload:     payload,
	}, nil
}

func (pss *PeerShareServer) SendResponsePacket(conn net.Conn, packet *PacketBody) error {
	buf := make([]byte, 0)
	buf = append(buf, byte(packet.MessageType))
	buf = append(buf, packet.Payload...)
	if len(buf) > utils.PacketMaxByteLength {
		return utils.ErrInvalidPacketSize
	}
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func (pss *PeerShareServer) SendResponse(conn net.Conn, msgType api.MessageType, msg any) error {
	respBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	packet := &PacketBody{
		MessageType: msgType,
		Payload:     respBytes,
	}
	if e := pss.SendResponsePacket(conn, packet); e != nil {
		return e
	}
	return nil
}
