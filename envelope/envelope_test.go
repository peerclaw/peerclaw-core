package envelope

import (
	"bytes"
	"testing"
	"time"

	"github.com/peerclaw/peerclaw-core/identity"
	"github.com/peerclaw/peerclaw-core/protocol"
)

func TestNewResponse_CopiesTraceID(t *testing.T) {
	req := New("alice", "bob", protocol.ProtocolA2A, []byte("hello"))
	resp := NewResponse(req, []byte("world"))

	if resp.TraceID != req.TraceID {
		t.Errorf("expected TraceID %s, got %s", req.TraceID, resp.TraceID)
	}
}

func TestNewResponse_SwapsSourceDest(t *testing.T) {
	req := New("alice", "bob", protocol.ProtocolA2A, []byte("hello"))
	resp := NewResponse(req, []byte("world"))

	if resp.Source != "bob" {
		t.Errorf("expected Source bob, got %s", resp.Source)
	}
	if resp.Destination != "alice" {
		t.Errorf("expected Destination alice, got %s", resp.Destination)
	}
}

func TestNewResponse_SetsMessageTypeResponse(t *testing.T) {
	req := New("alice", "bob", protocol.ProtocolA2A, []byte("hello"))
	resp := NewResponse(req, []byte("world"))

	if resp.MessageType != MessageTypeResponse {
		t.Errorf("expected MessageType response, got %s", resp.MessageType)
	}
}

func TestNewResponse_CopiesSessionID(t *testing.T) {
	req := New("alice", "bob", protocol.ProtocolA2A, []byte("hello"))
	req.SessionID = "session-123"
	resp := NewResponse(req, []byte("world"))

	if resp.SessionID != "session-123" {
		t.Errorf("expected SessionID session-123, got %s", resp.SessionID)
	}
}

func TestNewResponse_GeneratesNewID(t *testing.T) {
	req := New("alice", "bob", protocol.ProtocolA2A, []byte("hello"))
	resp := NewResponse(req, []byte("world"))

	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}
	if resp.ID == req.ID {
		t.Error("expected different ID from request")
	}
}

func TestNewResponse_CopiesProtocol(t *testing.T) {
	req := New("alice", "bob", protocol.ProtocolMCP, []byte("hello"))
	resp := NewResponse(req, []byte("world"))

	if resp.Protocol != protocol.ProtocolMCP {
		t.Errorf("expected Protocol mcp, got %s", resp.Protocol)
	}
}

func TestNewResponse_SetsPayload(t *testing.T) {
	req := New("alice", "bob", protocol.ProtocolA2A, []byte("hello"))
	resp := NewResponse(req, []byte("world"))

	if string(resp.Payload) != "world" {
		t.Errorf("expected Payload world, got %s", string(resp.Payload))
	}
}

func TestSigningPayload_Deterministic(t *testing.T) {
	env := &Envelope{
		Source:      "alice",
		Destination: "bob",
		Protocol:    protocol.ProtocolA2A,
		MessageType: MessageTypeRequest,
		Nonce:       "nonce-123",
		Timestamp:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Payload:     []byte("hello"),
	}

	p1 := env.SigningPayload()
	p2 := env.SigningPayload()
	if !bytes.Equal(p1, p2) {
		t.Error("SigningPayload should be deterministic")
	}
}

func TestSigningPayload_ChangesOnFieldMutation(t *testing.T) {
	env := &Envelope{
		Source:      "alice",
		Destination: "bob",
		Protocol:    protocol.ProtocolA2A,
		MessageType: MessageTypeRequest,
		Nonce:       "nonce-123",
		Timestamp:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Payload:     []byte("hello"),
	}

	baseline := env.SigningPayload()

	// Changing any covered field should produce a different signing payload.
	fields := []struct {
		name   string
		mutate func()
		revert func()
	}{
		{"Source", func() { env.Source = "mallory" }, func() { env.Source = "alice" }},
		{"Destination", func() { env.Destination = "eve" }, func() { env.Destination = "bob" }},
		{"Protocol", func() { env.Protocol = protocol.ProtocolMCP }, func() { env.Protocol = protocol.ProtocolA2A }},
		{"MessageType", func() { env.MessageType = MessageTypeResponse }, func() { env.MessageType = MessageTypeRequest }},
		{"Nonce", func() { env.Nonce = "other" }, func() { env.Nonce = "nonce-123" }},
		{"Timestamp", func() { env.Timestamp = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) }, func() { env.Timestamp = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }},
		{"Payload", func() { env.Payload = []byte("tampered") }, func() { env.Payload = []byte("hello") }},
	}

	for _, f := range fields {
		f.mutate()
		changed := env.SigningPayload()
		if bytes.Equal(baseline, changed) {
			t.Errorf("changing %s should produce different signing payload", f.name)
		}
		f.revert()
	}
}

func TestSignEnvelope_VerifyEnvelope_RoundTrip(t *testing.T) {
	kp, err := identity.GenerateKeypair()
	if err != nil {
		t.Fatalf("GenerateKeypair: %v", err)
	}

	env := New("alice", "bob", protocol.ProtocolA2A, []byte("hello"))
	env.Nonce = "test-nonce"
	identity.SignEnvelope(env, kp.PrivateKey)

	if env.Signature == "" {
		t.Fatal("expected non-empty signature")
	}

	pubKey, _ := identity.ParsePublicKey(kp.PublicKeyString())
	if err := identity.VerifyEnvelope(env, pubKey); err != nil {
		t.Fatalf("VerifyEnvelope should pass: %v", err)
	}

	// Tamper with Source — should fail.
	env.Source = "mallory"
	if err := identity.VerifyEnvelope(env, pubKey); err == nil {
		t.Error("tampered Source should fail verification")
	}
}
