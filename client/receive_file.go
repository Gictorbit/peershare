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

func (pc *PeerClient) ReceiveFile(code, outPath string) error {
	defer pc.Stop()
	pc.sharedCode = code
	if err := pc.InitPeerConnection(); err != nil {
		return err
	}

	// Register data channel creation handling
	pc.peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message, _ := utils.RandSeq(15)
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
				if err := pc.peerConnection.SetRemoteDescription(resp.Sdp); err != nil {
					log.Println(err)
					continue
				}
				log.Println("got offer")
				if e := pc.SendAnswer(); e != nil {
					log.Println("send answer failed", e)
					continue
				}
				pc.candidatesMux.Lock()
				for _, c := range pc.pendingCandidates {
					if signalCandidateErr := pc.SignalIceCandidate(c, pc.clientType); signalCandidateErr != nil {
						log.Println(signalCandidateErr)
					}
				}
				pc.candidatesMux.Unlock()
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
	err := pc.SendRequest(api.MessageTypeGetOfferRequest, &api.GetOfferRequest{Code: code})
	if err != nil {
		log.Fatal(err)
		return err
	}
	select {}
}
