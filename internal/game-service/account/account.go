package account

import "golang.org/x/crypto/bcrypt"

// Email represents an e-mail address.
type Email string

// Password is a secret value that is associated with an account and grants
// access to it and should therefore only be known by its owner.
type Password string

// PasswordMatcher searches for equality between two given Password's.
type PasswordMatcher func(p1, p2 Password) (bool, error)

// Account represents a user.
type Account struct {
	Email    Email
	Password Password
}

// Validate validates the Email string value. Returns whether
// the e-mail is a valid one or not.
func (email Email) Validate() bool {
	return true
}

// BasicPasswordMatcher is a PasswordMatcher that performs a simple
// value equality check to see if two Password values match up or not.
func BasicPasswordMatcher() PasswordMatcher {
	return func(p1, p2 Password) (b bool, e error) {
		return p1 == p2, nil
	}
}

// MatchPasswordsWithBCrypt uses the BCrypt algorithm to find equality
// in two Password inputs. May return an error, which is also an indicator
// that the two Password's are not a match.
func MatchPasswordsWithBCrypt() PasswordMatcher {
	return func(p1, p2 Password) (bool, error) {
		err := bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2))
		return err == nil, err
	}
}

// Validate validates the Password string value. Returns whether
// the password is a valid value or not.
func (password Password) Validate() bool {
	return len(password) >= 6
}
