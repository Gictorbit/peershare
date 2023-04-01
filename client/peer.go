package client

import (
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"google.golang.org/grpc"
)

type PeerClient struct {
	client pb.FileSharingServiceClient
}

func NewPeerClient(conn *grpc.ClientConn) *PeerClient {
	return &PeerClient{
		client: pb.NewFileSharingServiceClient(conn),
	}
}
