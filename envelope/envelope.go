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
	Signature   string            `json:"signature,omitempty"`
	Nonce       string            `json:"nonce,omitempty"`
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

// WithMetadata adds a metadata key-value pair and returns the envelope for chaining.
func (e *Envelope) WithMetadata(key, value string) *Envelope {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}
