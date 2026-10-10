package seeders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"

	authulamodels "github.com/Authula/authula/models"
	authulaservices "github.com/Authula/authula/services"
)

// Password bounds for the super user. The upper bound matches the
// email-password plugin, which would otherwise refuse the sign-in.
const (
	MinPasswordLen = 12
	MaxPasswordLen = 128
)

// ValidateSuperUser checks the SUPER_USER credentials before anything is
// seeded, because the account they create can sign in straight away.
func ValidateSuperUser(email, password string) error {
	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email {
		return errors.New("must be email:password with a valid email")
	}
	if len(password) < MinPasswordLen || len(password) > MaxPasswordLen {
		return fmt.Errorf("password must be %d to %d characters", MinPasswordLen, MaxPasswordLen)
	}
	return nil
}

// SuperUserSeeder creates the super user through Authula, so the row lands in
// the authula schema with a verified email and a hashed password, then
// mirrors it into the public users table.
type SuperUserSeeder struct {
	DB       *sql.DB
	Auth     *authulaservices.CoreServices
	Mirror   func(*authulamodels.User) error
	Email    string
	Password string
}

func (s SuperUserSeeder) Name() string {
	return "super user"
}

// Seed creates whatever is missing: the user, its email-password account and
// the mirror row. An existing account keeps its password.
func (s SuperUserSeeder) Seed() error {
	// The transaction only holds runSeed's advisory lock; Authula writes on
	// its own connection, so each step below is individually repeatable.
	return runSeed(s.DB, func(ctx context.Context, _ *sql.Tx) error {
		user, err := s.Auth.UserService.GetByEmail(ctx, s.Email)
		if err != nil {
			return err
		}
		if user == nil {
			if user, err = s.Auth.UserService.Create(ctx, "Super User", s.Email, true, nil, nil); err != nil {
				return err
			}
		}

		provider := authulamodels.AuthProviderEmail.String()
		account, err := s.Auth.AccountService.GetByUserIDAndProvider(ctx, user.ID, provider)
		if err != nil {
			return err
		}
		if account == nil {
			hash, err := s.Auth.PasswordService.Hash(s.Password)
			if err != nil {
				return err
			}
			if _, err := s.Auth.AccountService.Create(ctx, user.ID, user.Email, provider, &hash); err != nil {
				return err
			}
		}

		// Authula only logs a failed after-create hook, so mirror again here.
		return s.Mirror(user)
	})
}
