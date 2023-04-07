package client

import (
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"log"
	"net"
	"sync"
)

type PeerClient struct {
	listenAddr string
	conn       net.Conn
	wg         sync.WaitGroup
	sharedCode string
	clientType api.ClientType
}

type PeerShareClient interface {
	ReceiveFile(code, outPath string)
	SendFile(filePath string)
}

func NewPeerClient(listenAddr string, cliType api.ClientType) *PeerClient {
	return &PeerClient{
		listenAddr: listenAddr,
		wg:         sync.WaitGroup{},
		clientType: cliType,
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
	pc.conn.Close()
	log.Println("stop client...")
}
