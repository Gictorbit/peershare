package api

import (
	"github.com/pion/webrtc/v3"
)

type MessageType int32

const (
	MessagetypeMessageTypeSendOfferRequest   MessageType = 1
	MessagetypeMessageTypeSendOfferResponse  MessageType = 2
	MessagetypeMessageTypeGetOfferRequest    MessageType = 3
	MessagetypeMessageTypeGetOfferResponse   MessageType = 4
	MessagetypeMessageTypeSendAnswerRequest  MessageType = 5
	MessagetypeMessageTypeSendAnswerResponse MessageType = 6
)

type StatusCode int32

const (
	StatuscodeResponseCodeNotFound StatusCode = 0
	StatuscodeResponseCodeOk       StatusCode = 1
	StatuscodeResponseCodeError    StatusCode = 2
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
