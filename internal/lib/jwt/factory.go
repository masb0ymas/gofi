package jwt

import "gofi/internal/config"

func New(config *config.ConfigApp) *JWT {
	return &JWT{
		config: config,
	}
}
