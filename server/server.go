package server

import (
	"errors"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"go.uber.org/zap"
	"sync"
)

var (
	ErrInternalError = errors.New("")
)

type PeerShareService struct {
	pb.UnimplementedFileSharingServiceServer
	logger   *zap.Logger
	sessions map[string]*pb.SDP
	mu       sync.Mutex
}

func NewPeerShareService(logger *zap.Logger) *PeerShareService {
	return &PeerShareService{
		logger:   logger,
		sessions: make(map[string]*pb.SDP),
	}
}
