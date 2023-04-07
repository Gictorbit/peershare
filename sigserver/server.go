package sigserver

import (
	"encoding/json"
	"errors"
	api "github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/utils"
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
)

type Empty struct{}
type SessionPeers struct {
	Sender, Receiver *WebRTCPeer
}
type WebRTCPeer struct {
	Sdp           webrtc.SessionDescription
	IceCandidates []string
	Conn          net.Conn
}

type PeerShareServer struct {
	logger     *zap.Logger
	sessions   map[string]*SessionPeers
	mu         sync.Mutex
	listenAddr string
	ln         net.Listener
	quitChan   chan Empty
	wg         sync.WaitGroup
}

type PeerShareService interface {
	GetOfferHandler(req *api.GetOfferRequest, conn net.Conn) error
	SendAnswerHandler(req *api.SendAnswerRequest, conn net.Conn) error
	SendOfferHandler(req *api.SendOfferRequest, conn net.Conn) error
	SendCandidateHandler(req *api.SendIceCandidateRequest, conn net.Conn) error
}

var (
	_ PeerShareService = &PeerShareServer{}
)

func NewPeerShareServer(listenAddr string, logger *zap.Logger) *PeerShareServer {
	return &PeerShareServer{
		logger:     logger,
		sessions:   make(map[string]*SessionPeers),
		listenAddr: listenAddr,
		quitChan:   make(chan Empty),
		wg:         sync.WaitGroup{},
	}
}

func (pss *PeerShareServer) Start() {
	ln, err := net.Listen(utils.ServerSocketType, pss.listenAddr)
	if err != nil {
		pss.logger.Error("failed to listen", zap.Error(err))
		return
	}
	defer ln.Close()
	pss.ln = ln

	go pss.acceptConnections()
	pss.logger.Info("server started",
		zap.String("ListenAddress", pss.listenAddr),
	)
	<-pss.quitChan
}

func (pss *PeerShareServer) acceptConnections() {
	for {
		conn, err := pss.ln.Accept()
		if err != nil {
			pss.logger.Error("accept connection error", zap.Error(err))
			continue
		}
		pss.logger.Info("new Connection to the server", zap.String("Address", conn.RemoteAddr().String()))
		pss.wg.Add(1)
		go pss.HandleConnection(conn)
	}
}

func (pss *PeerShareServer) Stop() {
	pss.wg.Wait()
	pss.quitChan <- Empty{}
	pss.logger.Info("stop server")
}

func (pss *PeerShareServer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	defer pss.wg.Done()
	for {
		packet, err := pss.ReadPacket(conn)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				pss.logger.Error("read packet error", zap.Error(err))
			}
			return
		}
		switch packet.MessageType {
		case api.MessageTypeSendOfferRequest:
			req := &api.SendOfferRequest{}
			if e := json.Unmarshal(packet.Payload, req); e != nil {
				pss.logger.Error("unmarshal send offer request failed", zap.Error(e))
				continue
			}
			if e := pss.SendOfferHandler(req, conn); e != nil {
				pss.logger.Error("handle send offer failed", zap.Error(e))
				continue
			}
		case api.MessageTypeGetOfferRequest:
			req := &api.GetOfferRequest{}
			if e := json.Unmarshal(packet.Payload, req); e != nil {
				pss.logger.Error("unmarshal get offer request failed", zap.Error(e))
				continue
			}
			if e := pss.GetOfferHandler(req, conn); e != nil {
				pss.logger.Error("handle get offer failed", zap.Error(e))
				continue
			}
		case api.MessageTypeSendAnswerRequest:
			req := &api.SendAnswerRequest{}
			if e := json.Unmarshal(packet.Payload, req); e != nil {
				pss.logger.Error("unmarshal send answer request failed", zap.Error(e))
				continue
			}
			if e := pss.SendAnswerHandler(req, conn); e != nil {
				pss.logger.Error("handle send answer failed", zap.Error(e))
				continue
			}
		case api.MessageTypeSendIceCandidateRequest:
			req := &api.SendIceCandidateRequest{}
			if e := json.Unmarshal(packet.Payload, req); e != nil {
				pss.logger.Error("unmarshal ice candidate request failed", zap.Error(e))
				continue
			}
			if e := pss.SendCandidateHandler(req, conn); e != nil {
				pss.logger.Error("handle ice candidate failed", zap.Error(e))
				continue
			}
		}
	}
}
