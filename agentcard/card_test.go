package agentcard

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/peerclaw/peerclaw-core/protocol"
)

func TestHasCapability(t *testing.T) {
	c := &Card{Capabilities: []string{"llm", "tool-use"}}
	if !c.HasCapability("llm") {
		t.Error("expected HasCapability(llm) = true")
	}
	if c.HasCapability("unknown") {
		t.Error("expected HasCapability(unknown) = false")
	}
}

func TestSupportsProtocol(t *testing.T) {
	c := &Card{Protocols: []protocol.Protocol{protocol.ProtocolA2A, protocol.ProtocolMCP}}
	if !c.SupportsProtocol(protocol.ProtocolA2A) {
		t.Error("expected SupportsProtocol(a2a) = true")
	}
	if c.SupportsProtocol(protocol.ProtocolACP) {
		t.Error("expected SupportsProtocol(acp) = false")
	}
}

func TestHasSkill(t *testing.T) {
	c := &Card{
		Skills: []Skill{
			{Name: "summarize", Description: "Summarize text"},
			{Name: "translate", Description: "Translate text"},
		},
	}
	if !c.HasSkill("summarize") {
		t.Error("expected HasSkill(summarize) = true")
	}
	if c.HasSkill("code-review") {
		t.Error("expected HasSkill(code-review) = false")
	}
}

func TestHasTool(t *testing.T) {
	c := &Card{
		Tools: []Tool{
			{Name: "read_file", Description: "Read a file"},
			{Name: "search", Description: "Search the web"},
		},
	}
	if !c.HasTool("read_file") {
		t.Error("expected HasTool(read_file) = true")
	}
	if c.HasTool("write_file") {
		t.Error("expected HasTool(write_file) = false")
	}
}

func TestCardJSONRoundtrip(t *testing.T) {
	c := &Card{
		ID:           "test-id",
		Name:         "test-agent",
		Capabilities: []string{"llm"},
		Skills: []Skill{
			{Name: "summarize", InputModes: []string{"text"}, OutputModes: []string{"text"}},
		},
		Tools: []Tool{
			{Name: "search", Description: "web search", InputSchema: json.RawMessage(`{"type":"object"}`)},
		},
		Protocols: []protocol.Protocol{protocol.ProtocolA2A},
		Endpoint:  Endpoint{URL: "https://example.com"},
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Card
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Name != "test-agent" {
		t.Errorf("Name = %q", decoded.Name)
	}
	if len(decoded.Skills) != 1 || decoded.Skills[0].Name != "summarize" {
		t.Errorf("Skills = %+v", decoded.Skills)
	}
	if len(decoded.Tools) != 1 || decoded.Tools[0].Name != "search" {
		t.Errorf("Tools = %+v", decoded.Tools)
	}
}

func TestCardWithEmptySkillsTools(t *testing.T) {
	c := &Card{
		ID:       "test-id",
		Name:     "basic-agent",
		Endpoint: Endpoint{URL: "https://example.com"},
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}

	// Skills and Tools should be omitted from JSON when empty.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if _, exists := raw["skills"]; exists {
		t.Error("skills should be omitted when nil")
	}
	if _, exists := raw["tools"]; exists {
		t.Error("tools should be omitted when nil")
	}
}

func validPublicKey(t *testing.T) string {
	t.Helper()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(pub)
}

func TestValidateValid(t *testing.T) {
	c := &Card{
		Name:         "test-agent",
		PublicKey:    validPublicKey(t),
		Capabilities: []string{"llm"},
		Endpoint:     Endpoint{URL: "https://example.com", Port: 443},
		Protocols:    []protocol.Protocol{protocol.ProtocolA2A},
		Metadata:     map[string]string{"env": "prod"},
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid, got: %v", err)
	}
}

func TestValidateEmptyName(t *testing.T) {
	c := &Card{}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Errorf("expected name required error, got: %v", err)
	}
}

func TestValidateLongName(t *testing.T) {
	c := &Card{Name: strings.Repeat("x", 257)}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "at most 256") {
		t.Errorf("expected name length error, got: %v", err)
	}
}

func TestValidateControlCharsInName(t *testing.T) {
	c := &Card{Name: "test\x00agent"}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "control characters") {
		t.Errorf("expected control char error, got: %v", err)
	}
}

func TestValidateBadPublicKey(t *testing.T) {
	c := &Card{Name: "test", PublicKey: "not-valid-base64!!!"}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "valid base64") {
		t.Errorf("expected base64 error, got: %v", err)
	}

	// Wrong key size.
	c.PublicKey = base64.StdEncoding.EncodeToString([]byte("tooshort"))
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "32 bytes") {
		t.Errorf("expected key size error, got: %v", err)
	}
}

func TestValidateTooManyCapabilities(t *testing.T) {
	caps := make([]string, 51)
	for i := range caps {
		caps[i] = "cap"
	}
	c := &Card{Name: "test", Capabilities: caps}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "50 items") {
		t.Errorf("expected caps limit error, got: %v", err)
	}
}

func TestValidateLongCapability(t *testing.T) {
	c := &Card{Name: "test", Capabilities: []string{strings.Repeat("x", 129)}}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "128 characters") {
		t.Errorf("expected cap length error, got: %v", err)
	}
}

func TestValidateBadEndpointURL(t *testing.T) {
	c := &Card{Name: "test", Endpoint: Endpoint{URL: "ftp://bad"}}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "http/https") {
		t.Errorf("expected URL scheme error, got: %v", err)
	}
}

func TestValidateUnknownProtocol(t *testing.T) {
	c := &Card{Name: "test", Protocols: []protocol.Protocol{"grpc"}}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "unknown protocol") {
		t.Errorf("expected protocol error, got: %v", err)
	}
}

func TestValidateMetadataLimits(t *testing.T) {
	// Too many keys.
	meta := make(map[string]string)
	for i := 0; i < 51; i++ {
		meta[strings.Repeat("k", i+1)] = "v"
	}
	c := &Card{Name: "test", Metadata: meta}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "50 keys") {
		t.Errorf("expected metadata limit error, got: %v", err)
	}

	// Value too long.
	c2 := &Card{Name: "test", Metadata: map[string]string{"k": strings.Repeat("v", 1025)}}
	if err := c2.Validate(); err == nil || !strings.Contains(err.Error(), "1024") {
		t.Errorf("expected metadata value error, got: %v", err)
	}
}

func TestReputationConstants(t *testing.T) {
	if ReputationLow >= ReputationMedium {
		t.Error("ReputationLow should be less than ReputationMedium")
	}
	if ReputationMedium >= ReputationHigh {
		t.Error("ReputationMedium should be less than ReputationHigh")
	}
}
