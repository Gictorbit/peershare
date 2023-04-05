package api

import (
	"github.com/pion/webrtc/v3"
)

type MessageType int32

const (
	MessageTypeSendOfferRequest   MessageType = 1
	MessageTypeSendOfferResponse  MessageType = 2
	MessageTypeGetOfferRequest    MessageType = 3
	MessageTypeGetOfferResponse   MessageType = 4
	MessageTypeSendAnswerRequest  MessageType = 5
	MessageTypeSendAnswerResponse MessageType = 6
)

type StatusCode int32

const (
	ResponseCodeNotFound StatusCode = 0
	ResponseCodeOk       StatusCode = 1
	ResponseCodeError    StatusCode = 2
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
