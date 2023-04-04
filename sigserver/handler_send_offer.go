package sigserver

import (
	"crypto/rand"
	"encoding/base64"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) SendOfferHandler(req *pb.SendOfferRequest, conn net.Conn) error {
	code, err := pss.GenerateCode()
	if err != nil {
		pss.logger.Error("code generator failed", zap.Error(err))
		return pss.SendResponse(conn,
			pb.MessageType_MESSAGE_TYPE_SEND_OFFER_RESPONSE,
			&pb.SendOfferResponse{
				StatusCode: pb.StatusCode_RESPONSE_CODE_ERROR,
			})
	}
	pss.mu.Lock()
	defer pss.mu.Unlock()
	pss.sessions[code] = &PeerOffer{
		Sdp:  req.Sdp,
		Conn: conn,
	}
	return pss.SendResponse(conn,
		pb.MessageType_MESSAGE_TYPE_SEND_OFFER_RESPONSE,
		&pb.SendOfferResponse{
			Code:       code,
			StatusCode: pb.StatusCode_RESPONSE_CODE_OK,
		})
}

func (pss *PeerShareServer) GenerateCode() (string, error) {
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
