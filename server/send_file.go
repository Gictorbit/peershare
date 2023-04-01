package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pss *PeerShareService) SendFile(ctx context.Context, req *pb.SendFileRequest) (*pb.SendFileResponse, error) {
	code, err := pss.GenerateCode()
	if err != nil {
		pss.logger.Error("code generator failed", zap.Error(err))
		return nil, status.Error(codes.Internal, ErrInternalError.Error())
	}

	pss.mu.Lock()
	defer pss.mu.Unlock()
	pss.sessions[code] = req.Offer
	return &pb.SendFileResponse{
		Code: code,
	}, nil
}

func (pss *PeerShareService) GenerateCode() (string, error) {
	// Generate a random 15-byte string
	randomBytes := make([]byte, 15)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// Encode the random bytes as a base64 string
	code := base64.URLEncoding.EncodeToString(randomBytes)
	// Truncate the code to 10 characters
	return code, nil
}
