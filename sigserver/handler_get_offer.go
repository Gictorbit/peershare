package sigserver

import (
	"github.com/gictorbit/peershare/api"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) GetOfferHandler(req *api.GetOfferRequest, conn net.Conn) error {
	peers, found := pss.sessions[req.Code]
	if !found {
		pss.logger.Error("code not found", zap.String("code", req.Code))
		return pss.SendResponse(conn,
			api.MessageTypeGetOfferResponse,
			&api.GetOfferResponse{
				StatusCode: api.ResponseCodeError,
			})
	}
	peers.Receiver = &WebRTCPeer{
		Conn: conn,
	}
	return pss.SendResponse(conn,
		api.MessageTypeGetOfferResponse,
		&api.GetOfferResponse{
			Sdp:        peers.Sender.Sdp,
			StatusCode: api.ResponseCodeOk,
		})
}
