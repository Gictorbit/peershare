package client

import (
	"encoding/json"
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
	"log"
)

func (pc *PeerClient) ParseResponses(packet *PacketBody) error {
	switch packet.MessageType {
	case api.MessageTypeSendOfferResponse:
		resp := &api.SendOfferResponse{}
		if e := json.Unmarshal(packet.Payload, resp); e != nil || resp.StatusCode != api.ResponseCodeOk {
			return fmt.Errorf("unmarshal send offer response failed:%v\n", e)
		}
		fmt.Println("share code: ", resp.Code)
		pc.sharedCode = resp.Code
	case api.MessageTypeSendAnswerRequest:
		resp := &api.SendAnswerRequest{}
		if e := json.Unmarshal(packet.Payload, resp); e != nil {
			return fmt.Errorf("unmarshal transfer answer request failed:%v\n", e)
		}
		log.Println("got answer")
		if sdpErr := pc.peerConnection.SetRemoteDescription(resp.Sdp); sdpErr != nil {
			return fmt.Errorf("set answer to peer failed:%v", sdpErr)
		}
		pc.candidatesMux.Lock()
		for _, c := range pc.pendingCandidates {
			if signalCandidateErr := pc.SignalIceCandidate(c, pc.clientType); signalCandidateErr != nil {
				return signalCandidateErr
			}
		}
		pc.candidatesMux.Unlock()
	case api.MessageTypeGetOfferResponse:
		resp := &api.GetOfferResponse{}
		if e := json.Unmarshal(packet.Payload, resp); e != nil || resp.StatusCode != api.ResponseCodeOk {
			return fmt.Errorf("unmarshal get offer response failed:%v\n", e)

		}
		if err := pc.peerConnection.SetRemoteDescription(resp.Sdp); err != nil {
			return fmt.Errorf("set remote sdp failed:%v", err)
		}
		log.Println("got offer")
		if e := pc.SendAnswer(); e != nil {
			return fmt.Errorf("send answer failed:%v", e)

		}
		pc.candidatesMux.Lock()
		for _, c := range pc.pendingCandidates {
			if signalCandidateErr := pc.SignalIceCandidate(c, pc.clientType); signalCandidateErr != nil {
				return signalCandidateErr
			}
		}
		pc.candidatesMux.Unlock()
	case api.MessageTypeSendAnswerResponse:
		resp := &api.SendAnswerResponse{}
		if e := json.Unmarshal(packet.Payload, resp); e != nil {
			return fmt.Errorf("unmarshal send answer response failed:%v\n", e)
		}

	case api.MessageTypeTransferIceCandidate:
		resp := &api.TransferCandidates{}
		if e := json.Unmarshal(packet.Payload, resp); e != nil {
			return fmt.Errorf("unmarshal transfer candidate failed:%v\n", e)
		}
		for _, c := range resp.Candidates {
			candidate, err := utils.Decode(c)
			if err != nil {
				return fmt.Errorf("decode candidate failed:%v", err)
			}
			if addCandidErr := pc.peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate}); addCandidErr != nil {
				return fmt.Errorf("add candidate failed:%v", addCandidErr)
			}
		}
	case api.MessageTypeSendIceCandidateResponse:
		resp := &api.SendIceCandidateResponse{}
		if e := json.Unmarshal(packet.Payload, resp); e != nil || resp.StatusCode != api.ResponseCodeOk {
			return fmt.Errorf("unmarshal transfer candidate failed:%v", e)
		}
	default:
		return fmt.Errorf("not handled response")
	}
	return nil
}
