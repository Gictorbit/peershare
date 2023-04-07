package client

import (
	"fmt"
	logcharm "github.com/charmbracelet/log"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
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
	logger            *logcharm.Logger
}

type PeerShareClient interface {
	ReceiveFile(code, outPath string) error
	SendFile(filePath string) error
}

func NewPeerClient(listenAddr string, cliType api.ClientType, logger *logcharm.Logger) *PeerClient {
	return &PeerClient{
		listenAddr:        listenAddr,
		wg:                sync.WaitGroup{},
		clientType:        cliType,
		pendingCandidates: make([]*webrtc.ICECandidate, 0),
		doneChan:          make(chan Empty),
		logger:            logger,
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
		pc.logger.Error("close connection failed", "error", e)
	}
	pc.logger.Info("stop client...")
	if err := pc.peerConnection.Close(); err != nil {
		pc.logger.Error("cannot close peerConnection", "error", err)
	}
	os.Exit(0)
}
