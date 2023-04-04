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
			pb.StatusCode_RESPONSE_CODE_NOT_FOUND,
			&pb.ResponseError{
				Error: "code not found",
			})
	}
	return pss.SendResponse(conn,
		pb.MessageType_MESSAGE_TYPE_GET_OFFER_RESPONSE,
		pb.StatusCode_RESPONSE_CODE_OK,
		&pb.GetOfferResponse{
			Sdp: peerOffer.Sdp,
		})
}
