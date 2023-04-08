package client

import (
	"encoding/json"
	"errors"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
	"net"
	"os"
	"path/filepath"
	"strings"
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
						pc.logger.Error("error unmarshal file info", "error", err)
						return
					}
					pc.PrintFileInfo(fileInfo)
				} else {
					fPath := filepath.Join(outPath, fileInfo.Name)
					if err := os.WriteFile(fPath, msg.Data, 0644); err != nil {
						pc.logger.Error("error receive file", "error", err)
					}
					pc.logger.Info("file received successfully")
					info, err := utils.FileInfo(fPath)
					if err != nil {
						pc.logger.Error("get file info failed", "error", err)
						return
					}
					if info.Md5Sum == fileInfo.Md5Sum {
						pc.logger.Info("md5sum verified successfully")
						if e := d.SendText("success"); e != nil {
							pc.logger.Error("send text failed", "error", e)
						}
					} else {
						pc.logger.Error("md5sum not correct")
						if e := d.SendText("failed"); e != nil {
							pc.logger.Error("send text failed", "error", e, "md5sum", info.Md5Sum)
						}
					}
				}
			})
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
			//TODO fix invalid character error
			if e := pc.ParseResponses(packet); e != nil && !strings.Contains(e.Error(), "invalid character") {
				pc.logger.Error("parse response error", "error", e)
				continue
			}
		}
	}()
	err := pc.SendRequest(api.MessageTypeGetOfferRequest, &api.GetOfferRequest{Code: code})
	if err != nil {
		pc.logger.Error(err)
		return err
	}
	select {}
}
