package server

import (
	"context"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pss *PeerShareService) SendFile(ctx context.Context, req *pb.SendFileRequest) (*pb.SendFileResponse, error) {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		pss.logger.Error("random uuid failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	uniqueCode := randomUUID.String()[:10]
	pss.mu.Lock()
	defer pss.mu.Unlock()
	pss.sessions[uniqueCode] = req.Offer
	return &pb.SendFileResponse{
		Code: uniqueCode,
	}, nil
}
