package client

import (
	"encoding/json"
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/pion/webrtc/v3"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func (pc *PeerClient) ReceiveFile(code, outPath string) error {
	defer pc.Stop()
	pc.sharedCode = code
	if err := pc.InitPeerConnection(); err != nil {
		return err
	}

	// Register data channel creation handling
	pc.peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		if d.Label() == "file" {
			var (
				fileInfo   *api.File
				fInfoMutex sync.Mutex
			)
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				fInfoMutex.Lock()
				defer fInfoMutex.Unlock()
				if fileInfo == nil {
					fileInfo = &api.File{}
					err := json.Unmarshal(msg.Data, &fileInfo)
					if err != nil {
						log.Fatal("error unmarshal file info")
					}
					fmt.Println("fileInfo is:", fileInfo)
				} else {
					fPath := filepath.Join(outPath, fileInfo.Name)
					if err := os.WriteFile(fPath, msg.Data, 0644); err != nil {
						log.Fatal("error receive file", err)
					}
					fmt.Println("file received")
				}
			})
		}
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
