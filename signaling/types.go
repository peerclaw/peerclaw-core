package signaling

import "time"

// MessageType represents the type of signaling message.
type MessageType string

const (
	MessageTypeOffer        MessageType = "offer"
	MessageTypeAnswer       MessageType = "answer"
	MessageTypeICECandidate MessageType = "ice_candidate"
	MessageTypePing         MessageType = "ping"
	MessageTypePong         MessageType = "pong"
)

// SignalMessage is a signaling message exchanged between agents via the signaling server.
type SignalMessage struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	SDP       string      `json:"sdp,omitempty"`
	Candidate string      `json:"candidate,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}
