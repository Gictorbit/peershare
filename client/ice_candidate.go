package client

import (
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
)

func (pc *PeerClient) SignalIceCandidate(c *webrtc.ICECandidate, cliType api.ClientType) error {
	candidate, err := utils.Encode([]byte(c.ToJSON().Candidate))
	if err != nil {
		return fmt.Errorf("encode candidate failed:%v", err)
	}
	if sendCandidateErr := pc.SendRequest(api.MessageTypeSendIceCandidateRequest, &api.SendIceCandidateRequest{
		Candidate:  candidate,
		ClientType: cliType,
		Code:       pc.sharedCode,
	}); sendCandidateErr != nil {
		return fmt.Errorf("send ice candidate error:%v", sendCandidateErr)
	}
	return nil
}
