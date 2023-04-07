package client

import (
	"encoding/json"
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/pion/webrtc/v3"
	"log"
	"sync"
	"time"
)

func (pc *PeerClient) SendFile(filePath string) error {
	defer pc.conn.Close()
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
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
		}
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

		for range time.NewTicker(5 * time.Second).C {
			message := RandSeq(15)
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

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil || len(pc.sharedCode) == 0 {
			pendingCandidates = append(pendingCandidates, c)
		} else {
			if err := pc.SendRequest(api.MessageTypeSendIceCandidateRequest, &api.SendIceCandidateRequest{
				Candidate:  c.ToJSON().Candidate,
				ClientType: api.SenderClient,
				Code:       pc.sharedCode,
			}); err != nil {
				log.Println("send ice candidate error", err.Error())
			}
		}
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
				if sdpErr := peerConnection.SetRemoteDescription(resp.Sdp); sdpErr != nil {
					log.Println("set answer to remote desc", sdpErr)
					continue
				}
				candidatesMux.Lock()
				defer candidatesMux.Unlock()

				for _, c := range pendingCandidates {
					if err := pc.SendRequest(api.MessageTypeSendIceCandidateRequest, &api.SendIceCandidateRequest{
						Candidate:  c.ToJSON().Candidate,
						ClientType: api.SenderClient,
						Code:       pc.sharedCode,
					}); err != nil {
						log.Println("send ice candidate error", err.Error())
					}
				}
			case api.MessageTypeTransferIceCandidate:
				resp := &api.TransferCandidates{}
				if e := json.Unmarshal(packet.Payload, resp); e != nil {
					log.Printf("unmarshal transfer candidate failed:%v\n", e)
					continue
				}
				for _, candidate := range resp.Candidates {
					if addCandidErr := peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate}); addCandidErr != nil {
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
	// Create an offer to send to the other process
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}
	err = pc.SendRequest(api.MessageTypeSendOfferRequest, &api.SendOfferRequest{Sdp: offer})
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("sent offer to server")
	select {}
}
