package handlers

import (
	"strings"
	"testing"
)

// TestValidateURL_Accepts asserts the exact string validateURL hands to the
// store, because that string becomes the Location header on redirect.
func TestValidateURL_Accepts(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain host", "https://example.com", "https://example.com"},
		{"path and query", "http://example.com/a?b=1", "http://example.com/a?b=1"},
		{"port and subdomain", "https://sub.example.co.uk:8443/x", "https://sub.example.co.uk:8443/x"},
		{"surrounding whitespace is trimmed", "  https://example.com  ", "https://example.com"},
		{"anchor link", "https://example.com/docs#install", "https://example.com/docs#install"},
		{"hash-routed spa link", "https://app.example.com/#/dashboard", "https://app.example.com/#/dashboard"},
		{"fragment after query", "https://example.com/a?b=1#frag", "https://example.com/a?b=1#frag"},
		{"fragment containing a slash and dash", "https://example.com/p#a-b/c", "https://example.com/p#a-b/c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateURL(tt.in)
			if err != nil {
				t.Fatalf("validateURL(%q) returned error %v, want it accepted", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("validateURL(%q)\n got %q\nwant %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidateURL_Rejects(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr string
	}{
		{"no scheme", "example.com", "url must be absolute, like https://example.com"},
		{"protocol relative", "//evil.com", "url must be absolute, like https://example.com"},
		{"empty", "", "url must be absolute, like https://example.com"},
		{"only whitespace", "   ", "url must be absolute, like https://example.com"},
		{"javascript scheme", "javascript:alert(1)", "scheme should be http or https"},
		{"data scheme", "data:text/html,<script>", "scheme should be http or https"},
		{"ftp scheme", "ftp://example.com", "scheme should be http or https"},
		{"no host", "https://", "hostname should not be empty"},
		{"userinfo", "http://user:pw@example.com", "urls with username are not allowed"},
		{"too long", "https://example.com/" + strings.Repeat("a", maxURLLength), "url is too long"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateURL(tt.in)
			if err == nil {
				t.Fatalf("validateURL(%q) accepted and returned %q, want rejection", tt.in, got)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("validateURL(%q) error\n got %q\nwant %q", tt.in, err.Error(), tt.wantErr)
			}
		})
	}
}
