package client

import (
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
			if e := pc.ParseResponses(packet); e != nil {
				log.Println(err)
				continue
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
