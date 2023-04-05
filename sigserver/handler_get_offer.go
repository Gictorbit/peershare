package sigserver

import (
	"github.com/gictorbit/peershare/api"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) GetOfferHandler(req *api.GetOfferRequest, conn net.Conn) error {
	peerOffer, found := pss.sessions[req.Code]
	if !found {
		pss.logger.Error("code not found", zap.String("code", req.Code))
		return pss.SendResponse(conn,
			api.MessagetypeMessageTypeGetOfferResponse,
			&api.GetOfferResponse{
				StatusCode: api.StatuscodeResponseCodeError,
			})
	}
	return pss.SendResponse(conn,
		api.MessagetypeMessageTypeGetOfferResponse,
		&api.GetOfferResponse{
			Sdp:        peerOffer.Sdp,
			StatusCode: api.StatuscodeResponseCodeOk,
		})
}
