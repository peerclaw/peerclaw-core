package envelope

import (
	"time"

	"github.com/google/uuid"
	"github.com/peerclaw/peerclaw-core/protocol"
)

// MessageType categorizes the purpose of a message.
type MessageType string

const (
	MessageTypeRequest  MessageType = "request"
	MessageTypeResponse MessageType = "response"
	MessageTypeEvent    MessageType = "event"
	MessageTypeError    MessageType = "error"
)

// Envelope is the universal message wrapper that bridges different AI agent protocols.
// Every message flowing through PeerClaw is wrapped in an Envelope regardless of the
// originating protocol.
type Envelope struct {
	ID          string            `json:"id"`
	Source      string            `json:"source"`
	Destination string            `json:"destination"`
	Protocol    protocol.Protocol `json:"protocol"`
	MessageType MessageType       `json:"message_type"`
	ContentType string            `json:"content_type"`
	Payload     []byte            `json:"payload"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	TTL         int               `json:"ttl,omitempty"`
	TraceID     string            `json:"trace_id,omitempty"`
	SessionID    string            `json:"session_id,omitempty"`
	Signature    string            `json:"signature,omitempty"`
	Nonce        string            `json:"nonce,omitempty"`
	Encrypted    bool              `json:"encrypted,omitempty"`
	SenderX25519 string            `json:"sender_x25519,omitempty"`
}

// New creates a new Envelope with a generated ID and timestamp.
func New(source, destination string, proto protocol.Protocol, payload []byte) *Envelope {
	return &Envelope{
		ID:          uuid.New().String(),
		Source:      source,
		Destination: destination,
		Protocol:    proto,
		MessageType: MessageTypeRequest,
		ContentType: "application/json",
		Payload:     payload,
		Metadata:    make(map[string]string),
		Timestamp:   time.Now(),
		TraceID:     uuid.New().String(),
	}
}

// WithMessageType sets the message type and returns the envelope for chaining.
func (e *Envelope) WithMessageType(mt MessageType) *Envelope {
	e.MessageType = mt
	return e
}

// WithTTL sets the time-to-live and returns the envelope for chaining.
func (e *Envelope) WithTTL(ttl int) *Envelope {
	e.TTL = ttl
	return e
}

// WithSessionID sets the session ID for multi-turn conversations and returns the envelope for chaining.
func (e *Envelope) WithSessionID(sessionID string) *Envelope {
	e.SessionID = sessionID
	return e
}

// NewResponse creates a response envelope from a request envelope.
// It generates a new ID, copies TraceID and SessionID from the request,
// swaps Source and Destination, and sets MessageType to response.
func NewResponse(req *Envelope, payload []byte) *Envelope {
	return &Envelope{
		ID:          uuid.New().String(),
		Source:      req.Destination,
		Destination: req.Source,
		Protocol:    req.Protocol,
		MessageType: MessageTypeResponse,
		ContentType: req.ContentType,
		Payload:     payload,
		Metadata:    make(map[string]string),
		Timestamp:   time.Now(),
		TraceID:     req.TraceID,
		SessionID:   req.SessionID,
	}
}

// WithMetadata adds a metadata key-value pair and returns the envelope for chaining.
func (e *Envelope) WithMetadata(key, value string) *Envelope {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}
