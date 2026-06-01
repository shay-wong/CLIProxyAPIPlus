package management

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
)

func TestPutGeminiKeysPersistsNote(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &Handler{
		cfg:            &config.Config{},
		configFilePath: writeTestConfigFile(t),
	}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(
		http.MethodPut,
		"/v0/management/gemini-api-key",
		strings.NewReader(`[{"api-key":"gemini-key","note":" 中文备注 🚧 "}]`),
	)

	h.PutGeminiKeys(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := h.cfg.GeminiKey[0].Note; got != "中文备注 🚧" {
		t.Fatalf("gemini note = %q, want %q", got, "中文备注 🚧")
	}
}

func TestPatchProviderNotes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		path   string
		invoke func(*Handler, *gin.Context)
		cfg    *config.Config
		assert func(*testing.T, *config.Config)
	}{
		{
			name: "claude",
			path: "/v0/management/claude-api-key",
			invoke: func(h *Handler, c *gin.Context) {
				h.PatchClaudeKey(c)
			},
			cfg: &config.Config{ClaudeKey: []config.ClaudeKey{{APIKey: "claude-key", BaseURL: "https://claude.example.com", Headers: map[string]string{"X-Test": "1"}, Note: "old"}}},
			assert: func(t *testing.T, cfg *config.Config) {
				t.Helper()
				entry := cfg.ClaudeKey[0]
				if entry.Note != "新备注" {
					t.Fatalf("claude note = %q, want %q", entry.Note, "新备注")
				}
				if entry.BaseURL != "https://claude.example.com" || entry.Headers["X-Test"] != "1" {
					t.Fatalf("claude existing fields were not preserved: %+v", entry)
				}
			},
		},
		{
			name: "codex",
			path: "/v0/management/codex-api-key",
			invoke: func(h *Handler, c *gin.Context) {
				h.PatchCodexKey(c)
			},
			cfg: &config.Config{CodexKey: []config.CodexKey{{APIKey: "codex-key", BaseURL: "https://codex.example.com", Note: "old"}}},
			assert: func(t *testing.T, cfg *config.Config) {
				t.Helper()
				if got := cfg.CodexKey[0].Note; got != "新备注" {
					t.Fatalf("codex note = %q, want %q", got, "新备注")
				}
			},
		},
		{
			name: "vertex",
			path: "/v0/management/vertex-api-key",
			invoke: func(h *Handler, c *gin.Context) {
				h.PatchVertexCompatKey(c)
			},
			cfg: &config.Config{VertexCompatAPIKey: []config.VertexCompatKey{{APIKey: "vertex-key", BaseURL: "https://vertex.example.com", Note: "old"}}},
			assert: func(t *testing.T, cfg *config.Config) {
				t.Helper()
				if got := cfg.VertexCompatAPIKey[0].Note; got != "新备注" {
					t.Fatalf("vertex note = %q, want %q", got, "新备注")
				}
			},
		},
		{
			name: "openai compatibility",
			path: "/v0/management/openai-compatibility",
			invoke: func(h *Handler, c *gin.Context) {
				h.PatchOpenAICompat(c)
			},
			cfg: &config.Config{OpenAICompatibility: []config.OpenAICompatibility{{Name: "compat", BaseURL: "https://compat.example.com", Note: "old"}}},
			assert: func(t *testing.T, cfg *config.Config) {
				t.Helper()
				if got := cfg.OpenAICompatibility[0].Note; got != "新备注" {
					t.Fatalf("openai note = %q, want %q", got, "新备注")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				cfg:            tt.cfg,
				configFilePath: writeTestConfigFile(t),
			}
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(
				http.MethodPatch,
				tt.path,
				strings.NewReader(`{"index":0,"value":{"note":" 新备注 "}}`),
			)

			tt.invoke(h, c)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
			}
			tt.assert(t, h.cfg)
		})
	}
}

func TestOpenAICompatibilityWithAuthIndexIncludesNote(t *testing.T) {
	h := &Handler{
		cfg: &config.Config{OpenAICompatibility: []config.OpenAICompatibility{{
			Name:    "compat",
			BaseURL: "https://compat.example.com",
			Note:    "OpenAI 备注",
		}}},
	}

	entries := h.openAICompatibilityWithAuthIndex()
	payload, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("marshal response failed: %v", err)
	}

	if !strings.Contains(string(payload), `"note":"OpenAI 备注"`) {
		t.Fatalf("openai response payload does not include note: %s", payload)
	}
}
