package middlewares

import "testing"

func TestCrossOrigin(t *testing.T) {
	trusted := []string{"https://app.example.com"}
	cases := []struct {
		name                            string
		method, fetchSite, origin, host string
		want                            bool
	}{
		{"safe method", "GET", "cross-site", "https://evil.com", "api.example.com", false},
		{"non-browser client", "POST", "", "", "api.example.com", false},
		{"same origin", "POST", "same-origin", "https://api.example.com", "api.example.com", false},
		{"trusted origin", "POST", "same-site", "https://app.example.com", "api.example.com", false},
		{"sibling subdomain", "POST", "same-site", "https://blog.example.com", "api.example.com", true},
		{"cross site", "DELETE", "cross-site", "https://evil.com", "api.example.com", true},
		{"old browser, own origin", "POST", "", "https://api.example.com", "api.example.com", false},
		{"old browser, other origin", "PUT", "", "https://evil.com", "api.example.com", true},
	}
	for _, c := range cases {
		if got := crossOrigin(c.method, c.fetchSite, c.origin, c.host, trusted); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
