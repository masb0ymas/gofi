package argon2

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func (h *Argon2) Compare(encodedHash string, password string) (match bool, err error) {
	cfg, salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	generatedHash := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallel, cfg.KeyLength)
	return subtle.ConstantTimeCompare(generatedHash, hash) == 1, nil
}

func (h *Argon2) decodeHash(encodedHash string) (cfg argonConfig, salt []byte, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")

	if len(vals) != 6 {
		return argonConfig{}, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return argonConfig{}, nil, nil, ErrInvalidHash
	}

	if version != argon2.Version {
		return argonConfig{}, nil, nil, ErrInvalidHash
	}

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &cfg.Memory, &cfg.Iterations, &cfg.Parallel)
	if err != nil {
		return argonConfig{}, nil, nil, ErrInvalidHash
	}

	salt, err = base64.StdEncoding.DecodeString(vals[4])
	if err != nil {
		return argonConfig{}, nil, nil, ErrInvalidHash
	}
	cfg.SaltLength = uint32(len(salt))

	hash, err = base64.StdEncoding.DecodeString(vals[5])
	if err != nil {
		return argonConfig{}, nil, nil, ErrInvalidHash
	}
	cfg.KeyLength = uint32(len(hash))

	return cfg, salt, hash, nil
}
