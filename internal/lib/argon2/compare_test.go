package argon2

import (
	"encoding/base64"
	"fmt"
	"testing"

	"golang.org/x/crypto/argon2"
)

func TestDecodeHash(t *testing.T) {
	arg := &Argon2{}

	// Create a valid hash for testing
	salt := []byte("somesalt")
	key := argon2.IDKey([]byte("password"), salt, 1, 65536, 2, 32)
	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Key := base64.StdEncoding.EncodeToString(key)
	validHash := fmt.Sprintf("$argon2id$v=%d$m=65536,t=1,p=2$%s$%s", argon2.Version, b64Salt, b64Key)

	tests := []struct {
		name        string
		encodedHash string
		wantParams  *argonConfig
		wantErr     error
	}{
		{
			name:        "Valid hash",
			encodedHash: validHash,
			wantParams: &argonConfig{
				Memory:     65536,
				Iterations: 1,
				Parallel:   2,
				SaltLength: uint32(len(salt)),
				KeyLength:  uint32(len(key)),
			},
			wantErr: nil,
		},
		{
			name:        "Invalid format",
			encodedHash: "invalid",
			wantParams:  nil,
			wantErr:     ErrInvalidHash,
		},
		{
			name:        "Incompatible version",
			encodedHash: fmt.Sprintf("$argon2id$v=18$m=65536,t=1,p=2$%s$%s", b64Salt, b64Key),
			wantParams:  nil,
			wantErr:     ErrInvalidHash,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotParams, _, _, err := arg.decodeHash(tc.encodedHash)
			if err != tc.wantErr {
				t.Errorf("decodeHash() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if tc.wantParams != nil && gotParams != *tc.wantParams {
				t.Errorf("decodeHash() gotParams = %v, want %v", gotParams, tc.wantParams)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	arg := &Argon2{}

	// Helper function to create a valid hash
	createHash := func(password string) string {
		salt := []byte("somesalt")
		key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 2, 32)
		b64Salt := base64.StdEncoding.EncodeToString(salt)
		b64Key := base64.StdEncoding.EncodeToString(key)
		return fmt.Sprintf("$argon2id$v=%d$m=65536,t=1,p=2$%s$%s", argon2.Version, b64Salt, b64Key)
	}

	tests := []struct {
		name        string
		password    string
		encodedHash string
		wantMatch   bool
		wantErr     bool
	}{
		{
			name:        "Matching password",
			password:    "correct_password",
			encodedHash: createHash("correct_password"),
			wantMatch:   true,
			wantErr:     false,
		},
		{
			name:        "Non-matching password",
			password:    "incorrect_password",
			encodedHash: createHash("correct_password"),
			wantMatch:   false,
			wantErr:     false,
		},
		{
			name:        "Invalid hash",
			password:    "password",
			encodedHash: "invalid_hash",
			wantMatch:   false,
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMatch, err := arg.Compare(tc.password, tc.encodedHash)
			if (err != nil) != tc.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if gotMatch != tc.wantMatch {
				t.Errorf("Compare() gotMatch = %v, want %v", gotMatch, tc.wantMatch)
			}
		})
	}
}
