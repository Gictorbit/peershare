package client

import (
	"fmt"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
	"log"
	"os"
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
			os.Exit(0)
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

	err = pc.SendRequest(pb.MessageType_MESSAGE_TYPE_GET_OFFER_REQUEST, &pb.GetOfferRequest{
		Code: code,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	offerResp, err := utils.ReadMessageFromConn(pc.conn, &pb.GetOfferResponse{})
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
		SDP:  offerResp.Message.Sdp.Sdp,
		Type: webrtc.SDPType(offerResp.Message.Sdp.Type),
	}); err != nil {
		log.Fatal(err)
		return
	}
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}
	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}
	err = pc.SendRequest(pb.MessageType_MESSAGE_TYPE_SEND_ANSWER_REQUEST, &pb.SendAnswerRequest{
		Code: code,
		Sdp: &pb.SDP{
			Sdp:  answer.SDP,
			Type: uint32(answer.Type),
		},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	select {}
}
