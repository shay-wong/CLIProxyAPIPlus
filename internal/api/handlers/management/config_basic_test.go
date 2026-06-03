package management

import (
	"net/http"
	"strings"
	"testing"
)

func TestApplyGitHubHeadersUsesGitHubToken(t *testing.T) {
	t.Setenv(githubTokenEnvName, " test-token ")

	req, err := http.NewRequest(http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	applyGitHubHeaders(req, latestReleaseUserAgent)

	if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("Authorization = %q, want Bearer token", got)
	}
	if got := req.Header.Get("User-Agent"); got != latestReleaseUserAgent {
		t.Fatalf("User-Agent = %q, want %q", got, latestReleaseUserAgent)
	}
}

func TestLatestVersionStatusErrorMapsRateLimit(t *testing.T) {
	body := []byte(`{"message":"API rate limit exceeded for 127.0.0.1.","documentation_url":"https://docs.github.com/rest/overview/resources-in-the-rest-api#rate-limiting"}`)

	code, message := latestVersionStatusError(http.StatusForbidden, body)

	if code != "github_rate_limited" {
		t.Fatalf("code = %q, want github_rate_limited", code)
	}
	if !strings.Contains(message, "GitHub API rate limit exceeded") {
		t.Fatalf("message = %q, want friendly rate-limit message", message)
	}
	if strings.Contains(message, "unexpected_status") {
		t.Fatalf("message = %q, should not expose unexpected_status", message)
	}
}
