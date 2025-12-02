package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

type GoogleService struct {
	Config      *oauth2.Config
	RedisClient *redis.Client
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type AuthenticateResponse struct {
	Token    *oauth2.Token
	UserInfo *GoogleUserInfo
}

var (
	ErrGenerateState       = errors.New("error generate state token")
	ErrInvalidStateToken   = errors.New("invalid state token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

func (s GoogleService) generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error base64 rand: %s", err.Error())
	}

	state := base64.URLEncoding.EncodeToString(b)

	// Store state in Redis with 5 minute expiration
	ctx := context.Background()
	key := fmt.Sprintf("oauth:google:%s", state)

	err = s.RedisClient.Set(ctx, key, "1", 5*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("error storing state in redis: %s", err.Error())
	}

	return state, nil
}

func (s GoogleService) verifyStateToken(state string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("oauth:google:%s", state)

	// Check if state exists in Redis
	val, err := s.RedisClient.Get(ctx, key).Result()
	if err != nil || val == "" {
		return false
	}

	// Delete the state after verification (one-time use)
	s.RedisClient.Del(ctx, key)
	return true
}

func (s GoogleService) AuthCodeURL() (string, error) {
	state, err := s.generateStateToken()
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrGenerateState, err.Error())
	}

	// return url google oauth
	return s.Config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s GoogleService) Authenticate(state string, code string) (*AuthenticateResponse, error) {
	if !s.verifyStateToken(state) {
		return nil, ErrInvalidStateToken
	}

	token, err := s.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("error exchange code: %s", err.Error())
	}

	userInfo, err := s.getUserInfo(token)
	if err != nil {
		return nil, err
	}

	return &AuthenticateResponse{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}

func (s GoogleService) getUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("error get user info: %s", err.Error())
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("error decode user info: %s", err.Error())
	}

	return &userInfo, nil
}

func (s GoogleService) URLParse(rawURL string) (state string, code string, err error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("error parsing URL: %s", err.Error())
	}

	queryParams := parsedURL.Query()
	state = queryParams.Get("state")
	code = queryParams.Get("code")

	return state, code, nil
}

// RefreshAccessToken gets a new access token using a refresh token
// Returns a new token with updated AccessToken and Expiry
// The RefreshToken in the response may be the same or a new one (depending on provider)
func (s GoogleService) RefreshAccessToken(refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	// Create a token with only the refresh token
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	// Use TokenSource to automatically refresh the token
	tokenSource := s.Config.TokenSource(context.Background(), token)

	// Get the new token (this will use the refresh token to get a new access token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %s", err.Error())
	}

	// Verify that we got a valid access token
	if newToken.AccessToken == "" {
		return nil, fmt.Errorf("received empty access token")
	}

	return newToken, nil
}
