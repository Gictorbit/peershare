package client

import (
	"encoding/json"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"net"
)

type PacketBody struct {
	MessageType api.MessageType
	Payload     []byte
}

func (pc *PeerClient) ReadPacket(conn net.Conn) (*PacketBody, error) {
	buf := make([]byte, utils.PacketMaxByteLength)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return &PacketBody{
		MessageType: api.MessageType(buf[0]),
		Payload:     buf[1:n],
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

func (pc *PeerClient) SendRequest(msgType api.MessageType, msg any) error {
	respBytes, err := json.Marshal(msg)
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
