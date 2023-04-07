package sigserver

import (
	"github.com/gictorbit/peershare/api"
	"go.uber.org/zap"
	"net"
)

func (pss *PeerShareServer) SendCandidateHandler(req *api.SendIceCandidateRequest, conn net.Conn) error {
	peers, found := pss.sessions[req.Code]
	if !found {
		pss.logger.Error("code not found", zap.String("code", req.Code))
		return pss.SendResponse(conn,
			api.MessageTypeSendIceCandidateResponse,
			&api.SendIceCandidateResponse{
				StatusCode: api.ResponseCodeNotFound,
			})
	}
	if req.ClientType == api.SenderClient {
		if len(req.Candidate) > 0 {
			peers.Sender.IceCandidates = append(peers.Sender.IceCandidates, req.Candidate)
		}
		if peers.Receiver != nil {
			e := pss.SendResponse(peers.Receiver.Conn, api.MessageTypeTransferIceCandidate, &api.TransferCandidates{
				Candidates: peers.Sender.IceCandidates,
			})
			if e != nil {
				pss.logger.Error("transfer candidate error",
					zap.Error(e),
					zap.String("code", req.Code),
				)
				return pss.SendResponse(conn,
					api.MessageTypeSendIceCandidateResponse,
					&api.SendIceCandidateResponse{
						StatusCode: api.ResponseCodeError,
					})
			}
		}
	} else if req.ClientType == api.ReceiverClient {
		if len(req.Candidate) > 0 {
			peers.Receiver.IceCandidates = append(peers.Receiver.IceCandidates, req.Candidate)
		}
		if peers.Sender != nil {
			e := pss.SendResponse(peers.Sender.Conn, api.MessageTypeTransferIceCandidate, &api.TransferCandidates{
				Candidates: peers.Receiver.IceCandidates,
			})
			if e != nil {
				pss.logger.Error("transfer candidate error",
					zap.Error(e),
					zap.String("code", req.Code),
				)
				return pss.SendResponse(conn,
					api.MessageTypeSendIceCandidateResponse,
					&api.SendIceCandidateResponse{
						StatusCode: api.ResponseCodeError,
					})
			}
		}
	}
	return pss.SendResponse(conn,
		api.MessageTypeSendIceCandidateResponse,
		&api.SendIceCandidateResponse{
			StatusCode: api.ResponseCodeOk,
		})
}
