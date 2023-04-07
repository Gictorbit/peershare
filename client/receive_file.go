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

func (pc *PeerClient) ReceiveFile(code, outPath string) {
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
		if err := peerConnection.Close(); err != nil {
			fmt.Printf("cannot close peerConnection: %v\n", err)
		}
	}()

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

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message := RandSeq(15)
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendTextErr := d.SendText(message)
				if sendTextErr != nil {
					panic(sendTextErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
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
				ClientType: api.ReceiverClient,
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
			case api.MessageTypeGetOfferResponse:
				resp := &api.GetOfferResponse{}
				if e := json.Unmarshal(packet.Payload, resp); e != nil || resp.StatusCode != api.ResponseCodeOk {
					log.Printf("unmarshal get offer response failed:%v\n", e)
					continue
				}
				if err := peerConnection.SetRemoteDescription(resp.Sdp); err != nil {
					log.Println(err)
					continue
				}
				log.Println("got offer")
				answer, err := peerConnection.CreateAnswer(nil)
				if err != nil {
					log.Println(err)
					continue
				}
				// Sets the LocalDescription, and starts our UDP listeners
				err = peerConnection.SetLocalDescription(answer)
				if err != nil {
					log.Println(err)
					continue
				}
				err = pc.SendRequest(api.MessageTypeSendAnswerRequest, &api.SendAnswerRequest{
					Code: code,
					Sdp:  answer,
				})
				if err != nil {
					log.Fatal(err)
					return
				}
				candidatesMux.Lock()
				for _, c := range pendingCandidates {
					if err := pc.SendRequest(api.MessageTypeSendIceCandidateRequest, &api.SendIceCandidateRequest{
						Candidate:  c.ToJSON().Candidate,
						ClientType: api.SenderClient,
						Code:       pc.sharedCode,
					}); err != nil {
						log.Println("send ice candidate error", err.Error())
					}
				}
				candidatesMux.Unlock()
			case api.MessageTypeSendAnswerResponse:
				resp := &api.SendAnswerResponse{}
				if e := json.Unmarshal(packet.Payload, resp); e != nil {
					log.Printf("unmarshal send answer response failed:%v\n", e)
					continue
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
	err = pc.SendRequest(api.MessageTypeGetOfferRequest, &api.GetOfferRequest{Code: code})
	if err != nil {
		log.Fatal(err)
		return
	}
	select {}
}
