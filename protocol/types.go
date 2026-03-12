package protocol

// Protocol identifies an AI agent communication protocol.
type Protocol string

const (
	ProtocolA2A Protocol = "a2a" // Google Agent-to-Agent
	ProtocolACP Protocol = "acp" // IBM Agent Communication Protocol
	ProtocolMCP Protocol = "mcp" // Anthropic Model Context Protocol
)

// Transport identifies a network transport.
type Transport string

const (
	TransportHTTP  Transport = "http"
	TransportWS    Transport = "ws"
	TransportStdio Transport = "stdio"
)

// Valid returns true if the protocol is a known type.
func (p Protocol) Valid() bool {
	switch p {
	case ProtocolA2A, ProtocolACP, ProtocolMCP:
		return true
	}
	return false
}

func (p Protocol) String() string { return string(p) }

func (t Transport) String() string { return string(t) }
