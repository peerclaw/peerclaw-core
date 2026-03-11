package envelope

import (
	"testing"

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
