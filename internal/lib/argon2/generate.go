package argon2

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func (h *Argon2) Generate(password string) (string, error) {
	config := argonConfig{
		SaltLength: 16,
		KeyLength:  32,
		Iterations: 10,
		Memory:     64 * 1024,
		Parallel:   2,
	}

	encodedHash, err := h.generateFromPlainText(password, config)
	if err != nil {
		return "", err
	}

	return encodedHash, nil
}

func (h *Argon2) generateFromPlainText(password string, config argonConfig) (string, error) {
	salt, err := h.generateRandomBytes(config.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallel, config.KeyLength)

	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, config.Memory, config.Iterations, config.Parallel, b64Salt, b64Hash)

	return encodedHash, nil
}

func (h *Argon2) generateRandomBytes(length uint32) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
