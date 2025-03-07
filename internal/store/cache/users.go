package cache

import (
	"encoding/json"
	"fmt"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/bradfitz/gomemcache/memcache"
)

type UserStore struct {
	cacheDB *memcache.Client
}

func (s *UserStore) Get(ID int64) (*models.User, error) {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyUser, ID)

	item, err := s.cacheDB.Get(cacheKey)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}

		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Set(user *models.User) error {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyUser, user.ID)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	userItem := &memcache.Item{Key: cacheKey, Value: userBytes, Expiration: CacheTTLUser}
	if err := s.cacheDB.Set(userItem); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ID int64) error {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyUser, ID)

	err := s.cacheDB.Delete(cacheKey)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
