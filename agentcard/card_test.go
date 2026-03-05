package agentcard

import (
	"encoding/json"
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
