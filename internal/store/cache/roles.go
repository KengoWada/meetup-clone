package cache

import (
	"encoding/json"
	"fmt"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/bradfitz/gomemcache/memcache"
)

type RoleStore struct {
	cacheDB *memcache.Client
}

func (s *RoleStore) getCacheKey(ID int64) string {
	return fmt.Sprintf("%s:%d", CacheKeyRole, ID)
}

func (s *RoleStore) Get(ID int64) (*models.Role, error) {
	item, err := getFromCache(s.cacheDB, s.getCacheKey(ID))
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	var role models.Role
	if err := json.Unmarshal(item.Value, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleStore) Set(role *models.Role) error {
	roleBytes, err := json.Marshal(role)
	if err != nil {
		return err
	}

	roleItem := &memcache.Item{
		Key:        s.getCacheKey(role.ID),
		Value:      roleBytes,
		Expiration: CacheTTLOrg,
	}
	if err := s.cacheDB.Set(roleItem); err != nil {
		return err
	}

	return nil
}

func (s *RoleStore) Delete(ID int64) error {
	err := s.cacheDB.Delete(s.getCacheKey(ID))
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
