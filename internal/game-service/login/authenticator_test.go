package login

import (
	"testing"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
)

type dummyRepository struct {
	ReturnNil bool
}

func (repo *dummyRepository) Get(email account.Email, password account.Password) (*account.Account, error) {
	if repo.ReturnNil {
		return nil, nil
	} else {
		return &account.Account{Email: email, Password: password}, nil
	}
}

// Put puts the given Account under the specified Email into the Repository.
func (repo *dummyRepository) Put(email account.Email, account account.Account) error {
	return nil
}

func TestAuthenticator_Authenticat_Success(t *testing.T) {
	authenticator := NewAuthenticator(&dummyRepository{ReturnNil: false}, account.BasicPasswordMatcher())
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
	authenticator := NewAuthenticator(&dummyRepository{ReturnNil: true}, account.BasicPasswordMatcher())
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
	authenticator := NewAuthenticator(&dummyRepository{}, func(p1, p2 account.Password) (bool, error) {
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
