package managementasset

import "testing"

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
