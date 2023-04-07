package sigserver

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) SendOfferHandler(req *api.SendOfferRequest, conn net.Conn) error {
	code, err := utils.RandSeq(15)
	if err != nil {
		pss.logger.Error("code generator failed", zap.Error(err))
		return pss.SendResponse(conn,
			api.MessageTypeSendOfferResponse,
			&api.SendOfferResponse{
				StatusCode: api.ResponseCodeError,
			})
	}
	pss.mu.Lock()
	defer pss.mu.Unlock()
	pss.sessions[code] = &SessionPeers{
		Sender: &WebRTCPeer{
			Sdp:  req.Sdp,
			Conn: conn,
		},
	}
	return pss.SendResponse(conn,
		api.MessageTypeSendOfferResponse,
		&api.SendOfferResponse{
			Code:       code,
			StatusCode: api.ResponseCodeOk,
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
