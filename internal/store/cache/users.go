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

func (s *UserStore) getCacheKey(ID int64) string {
	return fmt.Sprintf("%s:%d", CacheKeyUser, ID)
}

func (s *UserStore) Get(ID int64) (*models.User, error) {
	item, err := getFromCache(s.cacheDB, s.getCacheKey(ID))
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	var user models.User
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Set(user *models.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	userItem := &memcache.Item{
		Key:        s.getCacheKey(user.ID),
		Value:      userBytes,
		Expiration: CacheTTLUser,
	}
	if err := s.cacheDB.Set(userItem); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ID int64) error {
	err := s.cacheDB.Delete(s.getCacheKey(ID))
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
