package character

import (
	"gitlab.com/pokesync/game-service/internal/game-service/account"
	"sync"
)

// Repository stores account associated player characters.
type Repository interface {
	Get(email account.Email) (*Profile, error)
	Put(email account.Email, profile *Profile) error
}

// InMemoryRepository is an in-memory implementation of a profile
// Repository where character profile records are forgotten about
// once the application's lifecycle ends.
type InMemoryRepository struct {
	profiles map[account.Email]*Profile
	mutex    *sync.Mutex
}

// NewInMemoryRepository constructs a new instance of an InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		profiles: make(map[account.Email]*Profile),
		mutex:    &sync.Mutex{},
	}
}

// Get looks up a Profile instance that may be registered under the specified
// Email. May return an error.
func (repo *InMemoryRepository) Get(email account.Email) (*Profile, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	profile := repo.profiles[email]
	if profile == nil {
		profile = &Profile{
			DisplayName: DisplayName(email),

			MapX:   0,
			MapZ:   2,
			LocalX: 60,
			LocalZ: 40,
		}
	}

	return profile, nil
}

// Put puts the given Profile under the specified Email into the Repository.
func (repo *InMemoryRepository) Put(email account.Email, profile *Profile) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.profiles[email] = profile
	return nil
}
