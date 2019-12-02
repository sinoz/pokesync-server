package login

import (
	"testing"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
)

type dummyFetcher struct {
	ReturnNil bool
}

func (repo *dummyFetcher) Get(email account.Email, password account.Password) (*account.Account, error) {
	if repo.ReturnNil {
		return nil, nil
	}

	return &account.Account{Email: email, Password: password}, nil
}

func TestAuthenticator_Authenticat_Success(t *testing.T) {
	authenticator := NewAuthenticator(&dummyFetcher{ReturnNil: false}, account.BasicPasswordMatcher())
	result, err := authenticator.Authenticate(account.Email("Sino@gmail.com"), account.Password("hello123"))
	if err != nil {
		t.Error(err)
	}

	switch result.(type) {
	case AuthSuccess:
		break
	default:
		t.Error("expected result to be of type AuthSuccess")
	}
}

func TestAuthenticator_Authenticat_CouldNotFindAccount(t *testing.T) {
	authenticator := NewAuthenticator(&dummyFetcher{ReturnNil: true}, account.BasicPasswordMatcher())
	result, err := authenticator.Authenticate(account.Email("Sino@gmail.com"), account.Password("hello123"))
	if err != nil {
		t.Error(err)
	}

	switch result.(type) {
	case CouldNotFindAccount:
		break
	default:
		t.Error("expected result to be of type CouldNotFindAccount")
	}
}

func TestAuthenticator_Authenticat_WrongPassword(t *testing.T) {
	authenticator := NewAuthenticator(&dummyFetcher{}, func(p1, p2 account.Password) (bool, error) {
		return false, nil
	})

	result, err := authenticator.Authenticate(account.Email("Sino@gmail.com"), account.Password("hello123"))
	if err != nil {
		t.Error(err)
	}

	switch result.(type) {
	case PasswordMismatch:
		break
	default:
		t.Error("expected result to be of type PasswordMismatch")
	}
}
