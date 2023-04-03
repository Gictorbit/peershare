package utils

import "errors"

const (
	PacketMaxByteLength = 2048
	ServerSocketType    = "tcp"
)

var (
	ErrInvalidPacketSize = errors.New("invalid packet size")
)
