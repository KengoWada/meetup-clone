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

func (s *OrganizationStore) getCacheKey(ID int64) string {
	return fmt.Sprintf("%s:%d", CacheKeyOrg, ID)
}

func (s *OrganizationStore) Get(ID int64) (*models.Organization, error) {
	item, err := getFromCache(s.cacheDB, s.getCacheKey(ID))
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
	orgBytes, err := json.Marshal(organization)
	if err != nil {
		return err
	}

	orgItem := &memcache.Item{
		Key:        s.getCacheKey(organization.ID),
		Value:      orgBytes,
		Expiration: CacheTTLOrg,
	}
	if err := s.cacheDB.Set(orgItem); err != nil {
		return err
	}

	return nil
}

func (s *OrganizationStore) Delete(ID int64) error {
	err := s.cacheDB.Delete(s.getCacheKey(ID))
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
