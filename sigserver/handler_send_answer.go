package sigserver

import (
	"github.com/gictorbit/peershare/api"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) SendAnswerHandler(req *api.SendAnswerRequest, conn net.Conn) error {
	peerSession, found := pss.sessions[req.Code]
	if !found {
		pss.logger.Error("code not found", zap.String("code", req.Code))
		return pss.SendResponse(conn,
			api.MessagetypeMessageTypeSendAnswerResponse,
			&api.SendAnswerResponse{
				StatusCode: api.StatuscodeResponseCodeNotFound,
			})
	}
	err := pss.SendResponse(peerSession.Conn,
		api.MessagetypeMessageTypeSendAnswerRequest, req)

	if err != nil {
		pss.logger.Error("failed to send answer to first peer",
			zap.Error(err),
			zap.String("code", req.Code),
		)
		return pss.SendResponse(conn,
			api.MessagetypeMessageTypeSendAnswerResponse,
			&api.SendAnswerResponse{
				StatusCode: api.StatuscodeResponseCodeError,
			})
	}
	return pss.SendResponse(conn,
		api.MessagetypeMessageTypeSendAnswerResponse,
		&api.SendAnswerResponse{
			StatusCode: api.StatuscodeResponseCodeOk,
		})
}
