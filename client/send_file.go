package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
	"net"
	"os"
)

type Empty struct{}

func (pc *PeerClient) SendFile(filePath string) error {
	defer pc.Stop()
	if err := pc.InitPeerConnection(); err != nil {
		return err
	}

	fileDataChannel, err := pc.peerConnection.CreateDataChannel("file", nil)
	if err != nil {
		return err
	}
	fileDataChannel.OnOpen(func() {
		if e := pc.SendFileToReceiver(fileDataChannel, filePath); e != nil {
			pc.logger.Error("SendFile Error", "error", e)
		}
	})
	defer fileDataChannel.Close()

	fileDataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		if string(msg.Data) == "success" {
			pc.logger.Info("sent file successfully")
			pc.Stop()
		}
	})
	go func() {
		for {
			packet, err := pc.ReadPacket(pc.conn)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					pc.logger.Error("error read packet", "error", err)
				}
				continue
			}
			if e := pc.ParseResponses(packet); e != nil {
				pc.logger.Error("parse response error", "error", e)
				continue
			}
		}
	}()
	if e := pc.SendNewOffer(); e != nil {
		return e
	}
	pc.logger.Info("sent offer to server")
	select {}
}

func (pc *PeerClient) SendFileToReceiver(dataChannel *webrtc.DataChannel, filePath string) error {
	fileInfo, err := utils.FileInfo(filePath)
	if err != nil {
		return err
	}
	pc.PrintFileInfo(fileInfo)
	bfileInfo, err := json.Marshal(fileInfo)
	if err != nil {
		return fmt.Errorf("marshal file info failed:%v", err)
	}
	if e := dataChannel.Send(bfileInfo); e != nil {
		return fmt.Errorf("failed to send file info:%v", e)
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file:%v", err)
	}
	if e := dataChannel.Send(file); e != nil {
		return fmt.Errorf("failed to send file:%v", e)
	}
	return nil
}
