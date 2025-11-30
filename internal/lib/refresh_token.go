package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gofi/internal/config"
)

type RefreshToken struct {
	config config.ConfigApp
}

func NewRefreshToken(config *config.ConfigApp) *RefreshToken {
	return &RefreshToken{
		config: *config,
	}
}

func (h RefreshToken) Generate(data string, expires int64) string {
	message := fmt.Sprintf("%s.%d", data, expires)

	hash := hmac.New(sha256.New, []byte(h.config.Secret))
	hash.Write([]byte(message))
	signature := hash.Sum(nil)

	encodedSignature := base64.URLEncoding.EncodeToString(signature)

	return encodedSignature
}

func (h RefreshToken) Verify(token string) (bool, error) {
	_, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false, fmt.Errorf("invalid signature encoding: %w", err)
	}

	return true, nil
}
