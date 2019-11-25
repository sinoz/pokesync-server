package account

import "sync"

// Repository stores every registered Account.
type Repository interface {
	Get(email Email, password Password) (*Account, error)
	Put(email Email, account Account) error
}

// InMemoryRepository is an in-memory implementation of an account
// Repository where account records are forgotten about once the
// application's lifecycle ends.
type InMemoryRepository struct {
	accounts map[Email]*Account
	mutex    *sync.Mutex
}

// NewInMemoryRepository constructs a new instance of an InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		accounts: make(map[Email]*Account),
		mutex:    &sync.Mutex{},
	}
}

// Get looks up an Account instance that may be registered under the specified
// Email. May return an error.
func (repo *InMemoryRepository) Get(email Email, password Password) (*Account, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	account := repo.accounts[email]
	if account == nil {
		account = &Account{
			Email:     email,
			Password:  "hello123",
			UserGroup: GameDesigner,
		}
	}

	return account, nil
}

// Put puts the given Account under the specified Email into the Repository.
func (repo *InMemoryRepository) Put(email Email, account Account) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.accounts[email] = &account
	return nil
}
