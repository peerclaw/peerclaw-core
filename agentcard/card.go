package agentcard

import (
	"encoding/base64"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
	"unicode"

	"github.com/peerclaw/peerclaw-core/protocol"
)

// AgentStatus represents the current status of an agent.
type AgentStatus string

const (
	StatusOnline   AgentStatus = "online"
	StatusBusy     AgentStatus = "busy"
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
	Skills        []Skill             `json:"skills,omitempty"`  // A2A-compatible structured skills
	Tools         []Tool              `json:"tools,omitempty"`   // MCP-compatible tool definitions
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
	NostrPubKey     string   `json:"nostr_pubkey,omitempty"`
	ReputationScore float64  `json:"reputation_score,omitempty"`
	NostrRelays     []string `json:"nostr_relays,omitempty"`
	InboxRelays     []string `json:"inbox_relays,omitempty"` // Nostr relay URLs for offline mailbox (NIP-65 style)
	IdentityAnchor  string   `json:"identity_anchor,omitempty"`
	PublicEndpoint  bool     `json:"public_endpoint,omitempty"` // Owner opts-in to expose endpoint URL publicly
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

// Skill represents a structured capability the agent can perform (A2A-compatible).
type Skill struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	InputModes  []string `json:"input_modes,omitempty"`
	OutputModes []string `json:"output_modes,omitempty"`
}

// Tool represents an MCP-compatible tool definition.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

// HasSkill checks if the agent advertises a given skill by name.
func (c *Card) HasSkill(name string) bool {
	for _, s := range c.Skills {
		if s.Name == name {
			return true
		}
	}
	return false
}

// HasTool checks if the agent advertises a given tool by name.
func (c *Card) HasTool(name string) bool {
	for _, t := range c.Tools {
		if t.Name == name {
			return true
		}
	}
	return false
}

// knownProtocols are the valid protocol values for agent registration.
var knownProtocols = map[string]bool{
	"a2a":      true,
	"mcp":      true,
	"acp":      true,
	"custom":   true,
	"peerclaw": true,
}

// Validate checks the Card fields for correctness and returns an error
// describing the first invalid field found.
func (c *Card) Validate() error {
	// Name: required, 1-256 chars, no control characters.
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(c.Name) > 256 {
		return fmt.Errorf("name must be at most 256 characters")
	}
	if containsControlChars(c.Name) {
		return fmt.Errorf("name must not contain control characters")
	}

	// PublicKey: if provided, must be valid base64-encoded Ed25519 key (32 bytes).
	if c.PublicKey != "" {
		keyBytes, err := base64.StdEncoding.DecodeString(c.PublicKey)
		if err != nil {
			return fmt.Errorf("public_key must be valid base64: %w", err)
		}
		if len(keyBytes) != ed25519.PublicKeySize {
			return fmt.Errorf("public_key must be %d bytes (Ed25519), got %d", ed25519.PublicKeySize, len(keyBytes))
		}
	}

	// Capabilities: max 50 items, each ≤128 chars.
	if len(c.Capabilities) > 50 {
		return fmt.Errorf("capabilities must have at most 50 items")
	}
	for _, cap := range c.Capabilities {
		if len(cap) > 128 {
			return fmt.Errorf("each capability must be at most 128 characters")
		}
	}

	// Endpoint URL: if provided, must be valid http/https URL.
	if c.Endpoint.URL != "" {
		u, err := url.Parse(c.Endpoint.URL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return fmt.Errorf("endpoint URL must be a valid http/https URL")
		}
	}

	// Endpoint Port: 0-65535.
	if c.Endpoint.Port < 0 || c.Endpoint.Port > 65535 {
		return fmt.Errorf("endpoint port must be between 0 and 65535")
	}

	// Protocols: max 10, must be known.
	if len(c.Protocols) > 10 {
		return fmt.Errorf("protocols must have at most 10 items")
	}
	for _, p := range c.Protocols {
		if !knownProtocols[string(p)] {
			return fmt.Errorf("unknown protocol: %s", p)
		}
	}

	// Metadata: max 50 keys, key ≤128, value ≤1024.
	if len(c.Metadata) > 50 {
		return fmt.Errorf("metadata must have at most 50 keys")
	}
	for k, v := range c.Metadata {
		if len(k) > 128 {
			return fmt.Errorf("metadata key must be at most 128 characters")
		}
		if len(v) > 1024 {
			return fmt.Errorf("metadata value must be at most 1024 characters")
		}
	}

	return nil
}

// containsControlChars returns true if s contains any Unicode control characters.
func containsControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
