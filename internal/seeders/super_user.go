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
// promotes it: Promote mirrors the user into the public users table and gives
// it the super user role.
type SuperUserSeeder struct {
	DB       *sql.DB
	Auth     *authulaservices.CoreServices
	Promote  func(*authulamodels.User) error
	Email    string
	Password string
}

func (s SuperUserSeeder) Name() string {
	return "super user"
}

// Seed creates whatever is missing: the user, its email-password account and
// the mirror row. SUPER_USER is the source of truth for the password: an
// existing account whose password differs is reset to it and its sessions are
// ended, so changing SUPER_USER locks out whoever holds a leaked password on
// the next boot.
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

		// Set when someone else may hold a session on this account: its
		// password was not SUPER_USER's, or it was registered unverified,
		// possibly by whoever guessed the email before the first seed.
		revoke := false

		// Password sign-in refuses unverified users; SUPER_USER vouches for this one.
		if !user.EmailVerified {
			user.EmailVerified = true
			if user, err = s.Auth.UserService.Update(ctx, user); err != nil {
				return err
			}
			revoke = true
		}

		provider := authulamodels.AuthProviderEmail.String()
		account, err := s.Auth.AccountService.GetByUserIDAndProvider(ctx, user.ID, provider)
		if err != nil {
			return err
		}
		if account == nil || account.Password == nil || !s.Auth.PasswordService.Verify(s.Password, *account.Password) {
			revoke = true
			hash, err := s.Auth.PasswordService.Hash(s.Password)
			if err != nil {
				return err
			}
			if account == nil {
				_, err = s.Auth.AccountService.Create(ctx, user.ID, user.Email, provider, &hash)
			} else {
				account.Password = &hash
				_, err = s.Auth.AccountService.Update(ctx, account)
			}
			if err != nil {
				return err
			}
		}

		// Ahead of the promotion: no session from before this seed may hold the role.
		if revoke {
			if err := s.Auth.SessionService.DeleteAllByUserID(ctx, user.ID); err != nil {
				return err
			}
		}

		return s.Promote(user)
	})
}
