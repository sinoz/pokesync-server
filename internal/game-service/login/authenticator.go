package login

import "gitlab.com/pokesync/game-service/internal/game-service/account"

// Authenticator authenticates users.
type Authenticator struct {
	AccountStorage  account.Repository
	PasswordMatcher account.PasswordMatcher
}

// couldNotFindAccount is a cached instance of the CouldNotFindAccount result.
var (
	couldNotFindAccount AuthResult = CouldNotFindAccount{}
	passwordMismatch    AuthResult = PasswordMismatch{}
)

// AuthSuccess is an AuthResult where the Authenticator has successfully
// authenticated a user with the provided e-mail/password combination.
type AuthSuccess struct {
	Account account.Account
}

// CouldNotFindAccount is an AuthResult of no Account record being associated
// with a provided Email address.
type CouldNotFindAccount struct{}

// CouldNotFindAccount is an AuthResult of an invalid password having
// been entered by the user.
type PasswordMismatch struct{}

// AuthResult is the result from attempting to authenticate a user.
type AuthResult interface{}

// NewAuthenticator constructs a new instance of an Authenticator.
func NewAuthenticator(repository account.Repository, matcher account.PasswordMatcher) Authenticator {
	return Authenticator{
		AccountStorage:  repository,
		PasswordMatcher: matcher,
	}
}

// Authenticate authenticates a user by the given E-mail / Password combination.
func (auth Authenticator) Authenticate(email account.Email, password account.Password) (AuthResult, error) {
	record, err := auth.AccountStorage.Get(email, password)
	if err != nil {
		return nil, err
	}

	if record == nil {
		return couldNotFindAccount, nil
	}

	matchFound, err := auth.PasswordMatcher(record.Password, password)
	if err != nil {
		return nil, err
	}

	if !matchFound {
		return passwordMismatch, nil
	}

	return AuthSuccess{Account: *record}, nil
}
