package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	authulamodels "github.com/Authula/authula/models"
	emailpasswordconstants "github.com/Authula/authula/plugins/email-password/constants"
)

func TestRefuseUnverified(t *testing.T) {
	// Same error as a wrong password, so the answer does not reveal the account.
	if err := refuseUnverified(context.Background(), &authulamodels.User{EmailVerified: false}); !errors.Is(err, emailpasswordconstants.ErrInvalidCredentials) {
		t.Errorf("unverified user: %v, want ErrInvalidCredentials", err)
	}
	if err := refuseUnverified(context.Background(), &authulamodels.User{EmailVerified: true}); err != nil {
		t.Errorf("verified user: %v", err)
	}
}

func TestPoolSplit(t *testing.T) {
	// The pools share the budget, and each keeps at least one connection.
	for n, want := range map[int][2]int{25: {13, 12}, 2: {1, 1}, 1: {1, 1}, 0: {1, 1}} {
		if api, auth := apiPoolSize(n), authPoolSize(n); api != want[0] || auth != want[1] {
			t.Errorf("split(%d) = api %d, auth %d; want %v", n, api, auth, want)
		}
	}
}

func TestAddTokensToResponse(t *testing.T) {
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"iat":1000,"exp":1900}`))
	reqCtx := &authulamodels.RequestContext{
		Values:       map[string]any{"access_token": "h." + payload + ".s", "refresh_token": "r"},
		ResponseBody: []byte(`{"user":{"id":"u1"}}`),
	}
	if err := addTokensToResponse(reqCtx); err != nil {
		t.Fatal(err)
	}

	var got struct {
		User         map[string]string `json:"user"`
		AccessToken  string            `json:"access_token"`
		RefreshToken string            `json:"refresh_token"`
		ExpiresAt    string            `json:"expires_at"`
		ExpiresIn    int64             `json:"expires_in"`
	}
	if err := json.Unmarshal(reqCtx.ResponseBody, &got); err != nil {
		t.Fatal(err)
	}
	if got.User["id"] != "u1" || got.AccessToken != "h."+payload+".s" || got.RefreshToken != "r" ||
		got.ExpiresAt != "1970-01-01T00:31:40Z" || got.ExpiresIn != 900 {
		t.Errorf("unexpected body: %s", reqCtx.ResponseBody)
	}

	// No token minted, or a redirect with no JSON body: left untouched.
	for _, rc := range []*authulamodels.RequestContext{
		{Values: map[string]any{}, ResponseBody: []byte(`{"a":1}`)},
		{Values: map[string]any{"access_token": "h." + payload + ".s"}},
	} {
		before := string(rc.ResponseBody)
		if err := addTokensToResponse(rc); err != nil || string(rc.ResponseBody) != before {
			t.Errorf("body changed to %q (err %v)", rc.ResponseBody, err)
		}
	}
}
