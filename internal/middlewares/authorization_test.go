package middlewares

import "testing"

func TestBearerToken(t *testing.T) {
	cases := map[string]string{
		"Bearer abc":   "abc",
		"bearer  abc ": "abc",
		"Bearer":       "",
		"Basic abc":    "",
		"abc":          "",
		"":             "",
	}
	for header, want := range cases {
		if got := bearerToken(header); got != want {
			t.Errorf("bearerToken(%q) = %q, want %q", header, got, want)
		}
	}
}
