package client

import (
	"fmt"
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
		go func() {
			for {
				packet, err := pc.ReadPacket(pc.conn)
				if err != nil {
					log.Printf("error read packet: %v", err)
					continue
				}
				if e := pc.ParseResponses(packet); e != nil {
					log.Println(err)
					continue
				}
			}
		}()
	}()
	if e := pc.SendNewOffer(); e != nil {
		return e
	}
	log.Println("sent offer to server")
	select {}
}
