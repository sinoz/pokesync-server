package login

import (
	"context"
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
)

// AccountProvider attempts to provide an Account by the given credentials.
type AccountProvider func(email account.Email, password account.Password) <-chan account.LoadResult

// AuthConfig holds configurations specific to the Authenticator.
type AuthConfig struct {
	AccountFetchTimeout time.Duration
}

// Authenticator authenticates users.
type Authenticator struct {
	Config          AuthConfig
	AccountProvider AccountProvider
	PasswordMatcher account.PasswordMatcher
}

var (
	couldNotFindAccount AuthResult = CouldNotFindAccount{}
	passwordMismatch    AuthResult = PasswordMismatch{}
	timedOut            AuthResult = TimedOut{}
)

// AuthSuccess is an AuthResult where the Authenticator has successfully
// authenticated a user with the provided e-mail/password combination.
type AuthSuccess struct {
	Account account.Account
}

// CouldNotFindAccount is an AuthResult of no Account record being associated
// with a provided Email address.
type CouldNotFindAccount struct {
	Error error
}

// PasswordMismatch is an AuthResult of an invalid password having
// been entered by the user.
type PasswordMismatch struct{}

// TimedOut is an AuthResult of the authentication procedure taking
// too long and has thus been 'timed out'.
type TimedOut struct{}

// AuthResult is the result from attempting to authenticate a user.
type AuthResult interface{}

// NewAuthenticator constructs a new instance of an Authenticator.
func NewAuthenticator(config AuthConfig, accountProvider AccountProvider, matcher account.PasswordMatcher) Authenticator {
	return Authenticator{
		Config:          config,
		AccountProvider: accountProvider,
		PasswordMatcher: matcher,
	}
}

// Authenticate authenticates a user by the given E-mail / Password combination.
func (auth Authenticator) Authenticate(ctx context.Context, email account.Email, password account.Password) (AuthResult, error) {
	select {
	case result := <-auth.AccountProvider(email, password):
		if result.Error != nil {
			return nil, result.Error
		}

		if result.Account == nil {
			return couldNotFindAccount, nil
		}

		matchFound, err := auth.PasswordMatcher(result.Account.Password, password)
		if err != nil {
			return nil, err
		}

		if !matchFound {
			return passwordMismatch, nil
		}

		return AuthSuccess{Account: *result.Account}, nil

	case <-ctx.Done():
		return timedOut, nil

	case <-time.After(auth.Config.AccountFetchTimeout):
		return timedOut, nil
	}
}
