package client

import (
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
)

func (pc *PeerClient) InitPeerConnection() error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		pc.logger.Info("Peer Connection State has changed", "state", s.String())
		if s == webrtc.PeerConnectionStateFailed {
			pc.logger.Warn("Peer Connection has gone to failed")
			pc.Stop()
		}
	})
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		pc.candidatesMux.Lock()
		defer pc.candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil || len(pc.sharedCode) == 0 {
			pc.pendingCandidates = append(pc.pendingCandidates, c)
		} else {
			if signalCandidateErr := pc.SignalIceCandidate(c, pc.clientType); signalCandidateErr != nil {
				pc.logger.Error(signalCandidateErr)
			}
		}
	})
	pc.peerConnection = peerConnection
	return nil
}

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

func (pc *PeerClient) SendNewOffer() error {
	// Create an offer to send to the other process
	offer, err := pc.peerConnection.CreateOffer(nil)
	if err != nil {
		return err
	}
	if err = pc.peerConnection.SetLocalDescription(offer); err != nil {
		return err
	}
	err = pc.SendRequest(api.MessageTypeSendOfferRequest, &api.SendOfferRequest{Sdp: offer})
	if err != nil {
		return err
	}
	return nil
}

func (pc *PeerClient) SendAnswer() error {
	answer, err := pc.peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}
	// Sets the LocalDescription, and starts our UDP listeners
	err = pc.peerConnection.SetLocalDescription(answer)
	if err != nil {
		return err
	}
	return pc.SendRequest(api.MessageTypeSendAnswerRequest, &api.SendAnswerRequest{
		Code: pc.sharedCode,
		Sdp:  answer,
	})
}
