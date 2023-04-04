package sigserver

import (
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) GetOfferHandler(req *pb.GetOfferRequest, conn net.Conn) error {
	peerOffer, found := pss.sessions[req.Code]
	if !found {
		pss.logger.Error("code not found", zap.String("code", req.GetCode()))
		return pss.SendResponse(conn,
			pb.MessageType_MESSAGE_TYPE_GET_OFFER_RESPONSE,
			&pb.GetOfferResponse{
				StatusCode: pb.StatusCode_RESPONSE_CODE_ERROR,
			})
	}
	return pss.SendResponse(conn,
		pb.MessageType_MESSAGE_TYPE_GET_OFFER_RESPONSE,
		&pb.GetOfferResponse{
			Sdp:        peerOffer.Sdp,
			StatusCode: pb.StatusCode_RESPONSE_CODE_OK,
		})
}
