package utils

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	api "github.com/gictorbit/peershare/api"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
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

func FileInfo(fPath string) (*api.File, error) {
	openedFile, err := os.Open(strings.TrimSpace(fPath)) // For read access.
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()
	fileExtension := filepath.Ext(fPath)
	fileStat, err := os.Stat(fPath)
	if err != nil {
		return nil, err
	}
	hash := md5.New()
	_, err = io.Copy(hash, openedFile)
	if err != nil {
		return nil, err
	}
	return &api.File{
		Name:      filepath.Base(openedFile.Name()),
		Size:      fileStat.Size(),
		Extension: fileExtension,
		Md5Sum:    fmt.Sprintf("%x", hash.Sum(nil)),
	}, nil
}
