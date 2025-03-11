package cache

import (
	"encoding/json"
	"fmt"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/bradfitz/gomemcache/memcache"
)

type OrganizationStore struct {
	cacheDB *memcache.Client
}

func (s *OrganizationStore) Get(ID int64) (*models.Organization, error) {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyOrg, ID)

	item, err := getFromCache(s.cacheDB, cacheKey)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	var organization models.Organization
	if err := json.Unmarshal(item.Value, &organization); err != nil {
		return nil, err
	}

	return &organization, nil
}

func (s *OrganizationStore) Set(organization *models.Organization) error {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyOrg, organization.ID)

	orgBytes, err := json.Marshal(organization)
	if err != nil {
		return err
	}

	orgItem := &memcache.Item{Key: cacheKey, Value: orgBytes, Expiration: CacheTTLOrg}
	if err := s.cacheDB.Set(orgItem); err != nil {
		return err
	}

	return nil
}

func (s *OrganizationStore) Delete(ID int64) error {
	cacheKey := fmt.Sprintf("%s:%d", CacheKeyOrg, ID)

	err := s.cacheDB.Delete(cacheKey)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
