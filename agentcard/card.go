package agentcard

import (
	"time"

	"github.com/peerclaw/peerclaw-go/protocol"
)

// AgentStatus represents the current status of an agent.
type AgentStatus string

const (
	StatusOnline   AgentStatus = "online"
	StatusOffline  AgentStatus = "offline"
	StatusDegraded AgentStatus = "degraded"
)

// Card represents an AI agent's identity and capabilities,
// compatible with the A2A Agent Card standard with PeerClaw extensions.
type Card struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Description   string              `json:"description,omitempty"`
	Version       string              `json:"version,omitempty"`
	PublicKey     string              `json:"public_key,omitempty"`
	Capabilities  []string            `json:"capabilities,omitempty"`
	Endpoint      Endpoint            `json:"endpoint"`
	Protocols     []protocol.Protocol `json:"protocols"`
	Auth          AuthInfo            `json:"auth,omitempty"`
	Metadata      map[string]string   `json:"metadata,omitempty"`
	PeerClaw      PeerClawExtension   `json:"peerclaw,omitempty"`
	Status        AgentStatus         `json:"status"`
	RegisteredAt  time.Time           `json:"registered_at"`
	LastHeartbeat time.Time           `json:"last_heartbeat"`
}

// Endpoint defines the network location for reaching an agent.
type Endpoint struct {
	URL       string             `json:"url"`
	Host      string             `json:"host,omitempty"`
	Port      int                `json:"port,omitempty"`
	Transport protocol.Transport `json:"transport,omitempty"`
}

// AuthInfo describes the authentication method for an agent.
type AuthInfo struct {
	Type   string            `json:"type,omitempty"` // bearer, mtls, api_key, none
	Params map[string]string `json:"params,omitempty"`
}

// PeerClawExtension contains PeerClaw-specific fields beyond the A2A standard.
type PeerClawExtension struct {
	NATType         string   `json:"nat_type,omitempty"`         // full_cone, restricted, symmetric, none
	RelayPreference string   `json:"relay_preference,omitempty"` // direct, relay, auto
	Priority        int      `json:"priority,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

// HasCapability checks if the agent advertises a given capability.
func (c *Card) HasCapability(cap string) bool {
	for _, cc := range c.Capabilities {
		if cc == cap {
			return true
		}
	}
	return false
}

// SupportsProtocol checks if the agent supports a given protocol.
func (c *Card) SupportsProtocol(p protocol.Protocol) bool {
	for _, pp := range c.Protocols {
		if pp == p {
			return true
		}
	}
	return false
}
