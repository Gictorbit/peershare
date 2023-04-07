package sigserver

import (
	"github.com/gictorbit/peershare/api"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) SendAnswerHandler(req *api.SendAnswerRequest, conn net.Conn) error {
	peers, found := pss.sessions[req.Code]
	if !found {
		pss.logger.Error("code not found", zap.String("code", req.Code))
		return pss.SendResponse(conn,
			api.MessageTypeSendAnswerResponse,
			&api.SendAnswerResponse{
				StatusCode: api.ResponseCodeNotFound,
			})
	}
	peers.Receiver.Sdp = req.Sdp
	err := pss.SendResponse(peers.Sender.Conn,
		api.MessageTypeSendAnswerRequest, req)

	if err != nil {
		pss.logger.Error("failed to send answer to first peer",
			zap.Error(err),
			zap.String("code", req.Code),
		)
		return pss.SendResponse(conn,
			api.MessageTypeSendAnswerResponse,
			&api.SendAnswerResponse{
				StatusCode: api.ResponseCodeError,
			})
	}
	return pss.SendResponse(conn,
		api.MessageTypeSendAnswerResponse,
		&api.SendAnswerResponse{
			StatusCode: api.ResponseCodeOk,
		})
}
