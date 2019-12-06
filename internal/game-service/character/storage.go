package character

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"gitlab.com/pokesync/game-service/internal/game-service/account"
)

// Repository stores account associated player characters.
type Repository interface {
	Get(email account.Email) (*Profile, error)
	Put(email account.Email, profile *Profile) error
}

// CacheConfig holds configurations specific to the Cache.
type CacheConfig struct {
	ExpireAfter time.Duration
}

// Cache temporarily stores character Profile's.
type Cache interface {
	Repository
}

// RedisCache is a type of Cache that stores character Profile's
// temporarily in a connected Redis instance.
type RedisCache struct {
	config      CacheConfig
	redisClient *redis.Client
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

// NewRedisCache constructs a new instance of a RedisCache.
func NewRedisCache(config CacheConfig, redisClient *redis.Client) Cache {
	return &RedisCache{
		config:      config,
		redisClient: redisClient,
	}
}

// Get attempts to fetch a character Profile that is stored under the
// the specified e-mail address.
func (repo *InMemoryRepository) Get(email account.Email) (*Profile, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	profile := repo.profiles[email]
	if profile == nil {
		profile = &Profile{
			DisplayName: DisplayName(email),

			Gender:    0,
			UserGroup: 5,

			MapX:   0,
			MapZ:   2,
			LocalX: 60,
			LocalZ: 40,
		}
	}

	return profile, nil
}

// Get attempts to fetch a character Profile that is stored under the
// specified e-mail address. May return an error if something went whilst
// trying to fetch the Profile from the cache.
func (cache *RedisCache) Get(email account.Email) (*Profile, error) {
	jsonString, err := cache.redisClient.
		Get(createRedisKey(email)).
		Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	profile := &Profile{}
	if err := json.Unmarshal([]byte(jsonString), profile); err != nil {
		return nil, err
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

// Put attempts to store the given character Profile into the Cache.
func (cache *RedisCache) Put(email account.Email, profile *Profile) error {
	profileBytes, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	_, err = cache.redisClient.
		Set(createRedisKey(email), string(profileBytes), cache.config.ExpireAfter).
		Result()

	return err
}

// createRedisKey creates a Redis key of the specified Email.
func createRedisKey(email account.Email) string {
	return fmt.Sprint("character-", string(email))
}
