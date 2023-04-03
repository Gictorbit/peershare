package server

import (
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) SendAnswerHandler(req *pb.SendAnswerRequest, conn net.Conn) error {
	peerSession, found := pss.sessions[req.GetCode()]
	if !found {
		pss.logger.Error("code not found", zap.String("code", req.GetCode()))
		return pss.SendResponse(conn,
			pb.MessageType_MESSAGE_TYPE_SEND_ANSWER_RESPONSE,
			pb.StatusCode_RESPONSE_CODE_NOT_FOUND,
			&pb.ResponseError{
				Error: "code not found",
			})
	}
	err := pss.SendResponse(peerSession.Conn,
		pb.MessageType_MESSAGE_TYPE_SEND_ANSWER_REQUEST,
		pb.StatusCode_RESPONSE_CODE_OK, req)

	if err != nil {
		pss.logger.Error("failed to send answer to first peer",
			zap.Error(err),
			zap.String("code", req.GetCode()),
		)
		return pss.SendResponse(conn,
			pb.MessageType_MESSAGE_TYPE_SEND_ANSWER_RESPONSE,
			pb.StatusCode_RESPONSE_CODE_ERROR,
			&pb.ResponseError{
				Error: "internal error",
			})
	}
	return pss.SendResponse(conn,
		pb.MessageType_MESSAGE_TYPE_SEND_ANSWER_RESPONSE,
		pb.StatusCode_RESPONSE_CODE_OK,
		&pb.SendAnswerResponse{})
}
