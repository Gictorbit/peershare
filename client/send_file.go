package client

import (
	"encoding/json"
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
	"log"
	"time"
)

func (pc *PeerClient) SendFile(filePath string) error {
	defer pc.Stop()
	if err := pc.InitPeerConnection(); err != nil {
		return err
	}

	// Create a datachannel with label 'data'
	dataChannel, err := pc.peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

		for range time.NewTicker(5 * time.Second).C {
			message, _ := utils.RandSeq(15)
			fmt.Printf("Sending '%s'\n", message)

			// Send the message as text
			sendTextErr := dataChannel.SendText(message)
			if sendTextErr != nil {
				panic(sendTextErr)
			}
		}
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	go func() {
		for {
			packet, err := pc.ReadPacket(pc.conn)
			if err != nil {
				log.Printf("error read packet: %v", err)
				continue
			}
			switch packet.MessageType {
			case api.MessageTypeSendOfferResponse:
				resp := &api.SendOfferResponse{}
				if e := json.Unmarshal(packet.Payload, resp); e != nil || resp.StatusCode != api.ResponseCodeOk {
					log.Printf("unmarshal send offer response failed:%v\n", e)
					continue
				}
				fmt.Println("share code: ", resp.Code)
				pc.sharedCode = resp.Code
			case api.MessageTypeSendAnswerRequest:
				resp := &api.SendAnswerRequest{}
				if e := json.Unmarshal(packet.Payload, resp); e != nil {
					log.Printf("unmarshal transfer answer request failed:%v\n", e)
					continue
				}
				log.Println("got answer")
				if sdpErr := pc.peerConnection.SetRemoteDescription(resp.Sdp); sdpErr != nil {
					log.Println("set answer to remote desc", sdpErr)
					continue
				}
				pc.candidatesMux.Lock()
				for _, c := range pc.pendingCandidates {
					if signalCandidateErr := pc.SignalIceCandidate(c, pc.clientType); signalCandidateErr != nil {
						log.Println(signalCandidateErr)
					}
				}
				pc.candidatesMux.Unlock()
			case api.MessageTypeTransferIceCandidate:
				resp := &api.TransferCandidates{}
				if e := json.Unmarshal(packet.Payload, resp); e != nil {
					log.Printf("unmarshal transfer candidate failed:%v\n", e)
					continue
				}

				for _, c := range resp.Candidates {
					candidate, err := utils.Decode(c)
					if err != nil {
						log.Println("error decode candidate", err)
					}
					if addCandidErr := pc.peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate}); addCandidErr != nil {
						log.Printf("add candidate failed:%v\n", addCandidErr)
						continue
					}
				}
			case api.MessageTypeSendIceCandidateResponse:
				resp := &api.SendIceCandidateResponse{}
				if e := json.Unmarshal(packet.Payload, resp); e != nil || resp.StatusCode != api.ResponseCodeOk {
					log.Printf("unmarshal transfer candidate failed:%v\n", e)
					continue
				}
			default:
				log.Println("not handled response")
			}
		}
	}()
	if e := pc.SendNewOffer(); e != nil {
		return e
	}
	log.Println("sent offer to server")
	select {}
}
