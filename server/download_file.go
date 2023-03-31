package server

import (
	"context"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pss *PeerShareService) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	sdp, found := pss.sessions[req.Code]
	if !found {
		return nil, status.Error(codes.NotFound, "code not found")
	}
	defer func() {
		pss.mu.Lock()
		delete(pss.sessions, req.Code)
		pss.mu.Unlock()
	}()
	return &pb.DownloadFileResponse{
		Offer: sdp,
	}, nil
}
