package argon2

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/argon2"
)

func TestGenerateRandomBytes(t *testing.T) {
	arg := &Argon2{}

	tests := []struct {
		name   string
		length uint32
	}{
		{"Generate 16 bytes", 16},
		{"Generate 32 bytes", 32},
		{"Generate 64 bytes", 64},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := arg.generateRandomBytes(tc.length)
			if err != nil {
				t.Errorf("generateRandomBytes() error = %v", err)
				return
			}
			if len(got) != int(tc.length) {
				t.Errorf("generateRandomBytes() returned %d bytes, want %d", len(got), tc.length)
			}
		})
	}
}

func TestGenerateFromPlainText(t *testing.T) {
	arg := &Argon2{}

	cfg := &argonConfig{
		SaltLength: 16,
		KeyLength:  32,
		Iterations: 10,
		Memory:     64 * 1024,
		Parallel:   2,
	}

	password := "testpassword"

	encodedHash, err := arg.generateFromPlainText(password, *cfg)
	if err != nil {
		t.Fatalf("generateFromPlainText() error = %v", err)
	}

	// Check the format of the encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		t.Errorf("Encoded hash has incorrect number of parts: got %d, want 6", len(parts))
	}

	if parts[1] != "argon2id" {
		t.Errorf("Incorrect algorithm identifier: got %s, want argon2id", parts[1])
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		t.Errorf("Error parsing version: %v", err)
	}
	if version != argon2.Version {
		t.Errorf("Incorrect version: got %d, want %d", version, argon2.Version)
	}

	var memory, iterations, parallelism int
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		t.Errorf("Error parsing parameters: %v", err)
	}
	if uint32(memory) != cfg.Memory || uint32(iterations) != cfg.Iterations || uint8(parallelism) != cfg.Parallel {
		t.Errorf("Incorrect parameters: got m=%d,t=%d,p=%d, want m=%d,t=%d,p=%d",
			memory, iterations, parallelism, cfg.Memory, cfg.Iterations, cfg.Parallel)
	}

	salt, err := base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		t.Errorf("Error decoding salt: %v", err)
	}
	if len(salt) != int(cfg.SaltLength) {
		t.Errorf("Incorrect salt length: got %d, want %d", len(salt), cfg.SaltLength)
	}

	hash, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		t.Errorf("Error decoding hash: %v", err)
	}
	if len(hash) != int(cfg.KeyLength) {
		t.Errorf("Incorrect hash length: got %d, want %d", len(hash), cfg.KeyLength)
	}
}

func TestGenerate(t *testing.T) {
	arg := &Argon2{}

	password := "testpassword"
	encodedHash, err := arg.Generate(password)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify that the generated hash can be successfully compared
	match, err := arg.Compare(password, encodedHash)
	if err != nil {
		t.Fatalf("Error comparing password: %v", err)
	}
	if !match {
		t.Errorf("Generated hash does not match original password")
	}

	// Verify that a different password doesn't match
	match, err = arg.Compare("wrongpassword", encodedHash)
	if err != nil {
		t.Fatalf("Error comparing password: %v", err)
	}
	if match {
		t.Errorf("Generated hash incorrectly matches wrong password")
	}
}
