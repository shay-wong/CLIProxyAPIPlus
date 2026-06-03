package managementasset

import "testing"

func TestGitHubBearerTokenPrefersGitHubToken(t *testing.T) {
	t.Setenv(githubTokenEnvName, " github-token ")
	t.Setenv("GITSTORE_GIT_URL", "https://github.com/example/repo.git")
	t.Setenv("GITSTORE_GIT_TOKEN", "gitstore-token")

	if got := githubBearerToken(); got != "github-token" {
		t.Fatalf("githubBearerToken() = %q, want github-token", got)
	}
}

func TestResolveReleaseURLForcesHTTPS(t *testing.T) {
	got := resolveReleaseURL("http://api.github.com/repos/example/panel")
	want := "https://api.github.com/repos/example/panel/releases/latest"
	if got != want {
		t.Fatalf("resolveReleaseURL() = %q, want %q", got, want)
	}
}

func TestManualRefreshUsesDistinctSingleflightKey(t *testing.T) {
	localPath := "/tmp/management.html"
	if got := managementAssetFlightKey(localPath, false); got != localPath {
		t.Fatalf("automatic key = %q, want %q", got, localPath)
	}
	if got := managementAssetFlightKey(localPath, true); got == localPath {
		t.Fatalf("forced key = %q, want distinct key", got)
	}
}

func TestResolveFallbackDownloadURL(t *testing.T) {
	tests := []struct {
		name string
		repo string
		want string
	}{
		{
			name: "empty repository uses default fork asset",
			repo: "",
			want: "https://github.com/shay-wong/Cli-Proxy-API-Management-Center/releases/latest/download/management.html",
		},
		{
			name: "github repository URL uses matching release asset",
			repo: "https://github.com/example/panel",
			want: "https://github.com/example/panel/releases/latest/download/management.html",
		},
		{
			name: "github git URL trims suffix",
			repo: "https://github.com/example/panel.git",
			want: "https://github.com/example/panel/releases/latest/download/management.html",
		},
		{
			name: "github API latest URL maps to release asset",
			repo: "https://api.github.com/repos/example/panel/releases/latest",
			want: "https://github.com/example/panel/releases/latest/download/management.html",
		},
		{
			name: "invalid repository falls back to default fork asset",
			repo: "://bad-url",
			want: "https://github.com/shay-wong/Cli-Proxy-API-Management-Center/releases/latest/download/management.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveFallbackDownloadURL(tt.repo); got != tt.want {
				t.Fatalf("resolveFallbackDownloadURL(%q) = %q, want %q", tt.repo, got, tt.want)
			}
		})
	}
}
