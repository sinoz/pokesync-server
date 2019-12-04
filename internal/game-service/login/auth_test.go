package login

import (
	"testing"
	"time"

	"gitlab.com/pokesync/game-service/internal/game-service/account"
)

type dummyFetcher struct {
	ReturnNil bool
}

func returnNilAccount(email account.Email, password account.Password) <-chan account.LoadResult {
	ch := make(chan account.LoadResult, 1)
	ch <- account.LoadResult{Account: nil, Error: nil}
	return ch
}

func returnMyAccount(email account.Email, password account.Password) <-chan account.LoadResult {
	ch := make(chan account.LoadResult, 1)
	ch <- account.LoadResult{Account: &account.Account{Email: email, Password: password}}
	return ch
}

func TestAuthenticator_Authenticate_Success(t *testing.T) {
	config := AuthConfig{AccountFetchTimeout: 1 * time.Second}
	authenticator := NewAuthenticator(config, returnMyAccount, account.BasicPasswordMatcher())
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
	config := AuthConfig{AccountFetchTimeout: 1 * time.Second}
	authenticator := NewAuthenticator(config, returnNilAccount, account.BasicPasswordMatcher())
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
	config := AuthConfig{AccountFetchTimeout: 1 * time.Second}
	authenticator := NewAuthenticator(config, returnMyAccount, func(p1, p2 account.Password) (bool, error) {
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
