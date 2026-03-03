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
	MessageTypeConfig       MessageType = "config"
)

// ICEServerConfig describes an ICE server (STUN or TURN) for WebRTC connectivity.
type ICEServerConfig struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

// SignalMessage is a signaling message exchanged between agents via the signaling server.
type SignalMessage struct {
	Type           MessageType      `json:"type"`
	From           string           `json:"from"`
	To             string           `json:"to"`
	SDP            string           `json:"sdp,omitempty"`
	Candidate      string           `json:"candidate,omitempty"`
	Timestamp      time.Time        `json:"timestamp"`
	ICEServers     []ICEServerConfig `json:"ice_servers,omitempty"`
	X25519PublicKey string           `json:"x25519_public_key,omitempty"`
}
