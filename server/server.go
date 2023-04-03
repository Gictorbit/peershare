package server

import (
	"errors"
	pb "github.com/gictorbit/peershare/api/gen/proto"
	"github.com/gictorbit/peershare/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
)

var (
	ErrInternalError = errors.New("internal error")
	ErrCodeNotFound  = errors.New("code not found")
)

type Empty struct{}

type PeerShareServer struct {
	logger     *zap.Logger
	sessions   map[string]*PeerOffer
	mu         sync.Mutex
	listenAddr string
	ln         net.Listener
	quitChan   chan Empty
	wg         sync.WaitGroup
}

type PeerOffer struct {
	Sdp  *pb.SDP
	Conn net.Conn
}

func NewPeerShareServer(listenAddr string, logger *zap.Logger) *PeerShareServer {
	return &PeerShareServer{
		logger:     logger,
		sessions:   make(map[string]*PeerOffer),
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
			pss.logger.Error("read packet error", zap.Error(err))
			return
		}
		switch packet.MessageType {
		case pb.MessageType_MESSAGE_TYPE_SEND_OFFER_REQUEST:
			req := &pb.SendOfferRequest{}
			if e := proto.Unmarshal(packet.Payload, req); e != nil {
				pss.logger.Error("unmarshal upload request failed", zap.Error(err))
				return
			}
			if e := pss.SendOfferHandler(req, conn); e != nil {
				pss.logger.Error("handle upload file failed", zap.Error(err))
				return
			}
		case pb.MessageType_MESSAGE_TYPE_GET_OFFER_REQUEST:
			req := &pb.GetOfferRequest{}
			if e := proto.Unmarshal(packet.Payload, req); e != nil {
				pss.logger.Error("unmarshal upload request failed", zap.Error(err))
				return
			}
			if e := pss.GetOfferHandler(req, conn); e != nil {
				pss.logger.Error("handle upload file failed", zap.Error(err))
				return
			}
		case pb.MessageType_MESSAGE_TYPE_SEND_ANSWER_REQUEST:
			req := &pb.SendAnswerRequest{}
			if e := proto.Unmarshal(packet.Payload, req); e != nil {
				pss.logger.Error("unmarshal upload request failed", zap.Error(err))
				return
			}
			if e := pss.SendAnswerHandler(req, conn); e != nil {
				pss.logger.Error("handle upload file failed", zap.Error(err))
				return
			}
		}
	}
}
