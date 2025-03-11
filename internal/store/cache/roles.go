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

func (s *RoleStore) Get(ID int64) (*models.Role, error) {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyRole, ID)

	item, err := getFromCache(s.cacheDB, cacheKey)
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
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyRole, role.ID)

	orgBytes, err := json.Marshal(role)
	if err != nil {
		return err
	}

	roleItem := &memcache.Item{Key: cacheKey, Value: orgBytes, Expiration: CacheTTLOrg}
	if err := s.cacheDB.Set(roleItem); err != nil {
		return err
	}

	return nil
}

func (s *RoleStore) Delete(ID int64) error {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyRole, ID)

	err := s.cacheDB.Delete(cacheKey)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
