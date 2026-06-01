package config

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestProviderNotesSurviveSerializationAndSanitize(t *testing.T) {
	cfg := &Config{
		GeminiKey: []GeminiKey{{APIKey: "gemini-key", Note: " Gemini 备注 🚧 "}},
		ClaudeKey: []ClaudeKey{{APIKey: "claude-key", Note: " Claude 备注 "}},
		CodexKey:  []CodexKey{{APIKey: "codex-key", BaseURL: "https://codex.example.com", Note: " Codex 备注 "}},
		OpenAICompatibility: []OpenAICompatibility{{
			Name:    "openai-compatible",
			BaseURL: "https://openai.example.com",
			Note:    " OpenAI 备注 ",
		}},
		VertexCompatAPIKey: []VertexCompatKey{{APIKey: "vertex-key", Note: " Vertex 备注 "}},
	}

	cfg.SanitizeGeminiKeys()
	cfg.SanitizeClaudeKeys()
	cfg.SanitizeCodexKeys()
	cfg.SanitizeOpenAICompatibility()
	cfg.SanitizeVertexCompatKeys()

	if got := cfg.GeminiKey[0].Note; got != "Gemini 备注 🚧" {
		t.Fatalf("gemini note = %q, want %q", got, "Gemini 备注 🚧")
	}
	if got := cfg.ClaudeKey[0].Note; got != "Claude 备注" {
		t.Fatalf("claude note = %q, want %q", got, "Claude 备注")
	}
	if got := cfg.CodexKey[0].Note; got != "Codex 备注" {
		t.Fatalf("codex note = %q, want %q", got, "Codex 备注")
	}
	if got := cfg.OpenAICompatibility[0].Note; got != "OpenAI 备注" {
		t.Fatalf("openai note = %q, want %q", got, "OpenAI 备注")
	}
	if got := cfg.VertexCompatAPIKey[0].Note; got != "Vertex 备注" {
		t.Fatalf("vertex note = %q, want %q", got, "Vertex 备注")
	}

	yamlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("yaml marshal failed: %v", err)
	}
	var yamlRoundTrip Config
	if err := yaml.Unmarshal(yamlBytes, &yamlRoundTrip); err != nil {
		t.Fatalf("yaml unmarshal failed: %v", err)
	}
	if got := yamlRoundTrip.GeminiKey[0].Note; got != "Gemini 备注 🚧" {
		t.Fatalf("yaml round-trip gemini note = %q", got)
	}
	if got := yamlRoundTrip.OpenAICompatibility[0].Note; got != "OpenAI 备注" {
		t.Fatalf("yaml round-trip openai note = %q", got)
	}

	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}
	var jsonRoundTrip Config
	if err := json.Unmarshal(jsonBytes, &jsonRoundTrip); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	if got := jsonRoundTrip.CodexKey[0].Note; got != "Codex 备注" {
		t.Fatalf("json round-trip codex note = %q", got)
	}
	if got := jsonRoundTrip.VertexCompatAPIKey[0].Note; got != "Vertex 备注" {
		t.Fatalf("json round-trip vertex note = %q", got)
	}
}
