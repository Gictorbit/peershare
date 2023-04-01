package main

import (
	"context"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"github.com/gictorbit/peershare/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime/debug"
)

func RunSignalingGRPCServer(grpcAddr string) *grpc.Server {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sigServer := server.NewPeerShareService(logger)

	grpcServer := NewGrpcServer(grpcAddr, logger)
	pb.RegisterFileSharingServiceServer(grpcServer, sigServer) // register your service implementation

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Printf("starting user gRPC server on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return grpcServer
}

func NewGrpcServer(grpcAddr string, logger *zap.Logger) *grpc.Server {
	unaryServerOptions := []grpc.UnaryServerInterceptor{
		grpcCtxTags.UnaryServerInterceptor(grpcCtxTags.WithFieldExtractor(grpcCtxTags.CodeGenRequestFieldExtractor)),
		grpcZap.UnaryServerInterceptor(logger),
		grpcZap.PayloadUnaryServerInterceptor(logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return LogRequest == true
		}),
		grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("stack trace from panic " + string(debug.Stack()))
			return status.Errorf(codes.Internal, "%v", p)
		})),
	}
	streamServerOptions := []grpc.StreamServerInterceptor{
		grpcCtxTags.StreamServerInterceptor(grpcCtxTags.WithFieldExtractor(grpcCtxTags.CodeGenRequestFieldExtractor)),
		grpcZap.StreamServerInterceptor(logger),
		grpcZap.PayloadStreamServerInterceptor(logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return LogRequest == true
		}),
		grpcRecovery.StreamServerInterceptor(grpcRecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("stack trace from panic " + string(debug.Stack()))
			return status.Errorf(codes.Internal, "%v", p)
		})),
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamServerOptions...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryServerOptions...)),
	)
	return grpcServer
}
