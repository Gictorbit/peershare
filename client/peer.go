package client

import (
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"google.golang.org/grpc"
)

type PeerClient struct {
	client pb.FileSharingServiceClient
}

type PeerShareClient interface {
	ReceiveFile(code, outPath string)
	SendFile(filePath string)
}

func NewPeerClient(conn *grpc.ClientConn) *PeerClient {
	return &PeerClient{
		client: pb.NewFileSharingServiceClient(conn),
	}
}
