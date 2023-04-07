package client

import (
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
	"log"
	"net"
	"os"
	"sync"
)

type PeerClient struct {
	listenAddr        string
	conn              net.Conn
	wg                sync.WaitGroup
	sharedCode        string
	clientType        api.ClientType
	candidatesMux     sync.Mutex
	pendingCandidates []*webrtc.ICECandidate
	peerConnection    *webrtc.PeerConnection
	doneChan          chan Empty
}

type PeerShareClient interface {
	ReceiveFile(code, outPath string) error
	SendFile(filePath string) error
}

func NewPeerClient(listenAddr string, cliType api.ClientType) *PeerClient {
	return &PeerClient{
		listenAddr:        listenAddr,
		wg:                sync.WaitGroup{},
		clientType:        cliType,
		pendingCandidates: make([]*webrtc.ICECandidate, 0),
		doneChan:          make(chan Empty),
	}
}

func (pc *PeerClient) Connect() error {
	conn, err := net.Dial(utils.ServerSocketType, pc.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to dial server: %v\n", err.Error())
	}
	pc.conn = conn
	return nil
}

func (pc *PeerClient) Stop() {
	pc.wg.Wait()
	if e := pc.conn.Close(); e != nil {
		log.Printf("close connection failed:%v\n", e)
	}
	log.Println("stop client...")
	if err := pc.peerConnection.Close(); err != nil {
		fmt.Printf("cannot close peerConnection: %v\n", err)
	}
	os.Exit(0)
}
