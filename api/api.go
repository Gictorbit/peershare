package api

import (
	"github.com/pion/webrtc/v3"
)

type MessageType int32

const (
	MessageTypeSendOfferRequest         MessageType = 1
	MessageTypeSendOfferResponse        MessageType = 2
	MessageTypeGetOfferRequest          MessageType = 3
	MessageTypeGetOfferResponse         MessageType = 4
	MessageTypeSendAnswerRequest        MessageType = 5
	MessageTypeSendAnswerResponse       MessageType = 6
	MessageTypeSendIceCandidateRequest  MessageType = 7
	MessageTypeSendIceCandidateResponse MessageType = 8
	MessageTypeTransferIceCandidate     MessageType = 9
	MessageTypeTransferAnswer           MessageType = 10
)

type StatusCode int32

const (
	ResponseCodeNotFound StatusCode = 0
	ResponseCodeOk       StatusCode = 1
	ResponseCodeError    StatusCode = 2
)

type ClientType int32

const (
	SenderClient   ClientType = 0
	ReceiverClient ClientType = 1
)

type SendOfferRequest struct {
	Sdp webrtc.SessionDescription `json:"sdp,omitempty"`
}

type SendOfferResponse struct {
	Code       string     `json:"code,omitempty"`
	StatusCode StatusCode `json:"status_code,omitempty"`
}

type GetOfferRequest struct {
	Code string `json:"code,omitempty"`
}

type GetOfferResponse struct {
	Sdp        webrtc.SessionDescription `json:"sdp,omitempty"`
	StatusCode StatusCode                `json:"status_code,omitempty"`
}

type SendAnswerRequest struct {
	Code string                    `json:"code,omitempty"`
	Sdp  webrtc.SessionDescription `json:"sdp,omitempty"`
}

type SendAnswerResponse struct {
	StatusCode StatusCode `json:"status_code,omitempty"`
}

type SendIceCandidateRequest struct {
	Code       string     `json:"code,omitempty"`
	Candidate  string     `json:"candidate,omitempty"`
	ClientType ClientType `json:"client_type,omitempty"`
}

type SendIceCandidateResponse struct {
	StatusCode StatusCode `json:"status_code,omitempty"`
}

type TransferCandidates struct {
	Candidates []string `json:"candidates,omitempty"`
}
