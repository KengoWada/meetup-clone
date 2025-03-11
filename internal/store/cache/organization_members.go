package cache

import (
	"encoding/json"
	"fmt"

	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/bradfitz/gomemcache/memcache"
)

type OrganizationMemberStore struct {
	cacheDB *memcache.Client
}

func (s *OrganizationMemberStore) getCacheKey(userID, orgID int64) string {
	return fmt.Sprintf("%s:%d,%d", CacheKeyOrgMember, userID, orgID)
}

func (s *OrganizationMemberStore) Get(userID, orgID int64) (*models.OrganizationMember, error) {
	item, err := getFromCache(s.cacheDB, s.getCacheKey(userID, orgID))
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	var member models.OrganizationMember
	if err := json.Unmarshal(item.Value, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

func (s *OrganizationMemberStore) Set(member *models.OrganizationMember) error {
	orgMemberBytes, err := json.Marshal(member)
	if err != nil {
		return err
	}

	roleItem := &memcache.Item{
		Key:        s.getCacheKey(member.UserProfileID, member.OrganizationID),
		Value:      orgMemberBytes,
		Expiration: CacheTTLOrg,
	}
	if err := s.cacheDB.Set(roleItem); err != nil {
		return err
	}

	return nil
}

func (s *OrganizationMemberStore) Delete(userID, orgID int64) error {
	err := s.cacheDB.Delete(s.getCacheKey(userID, orgID))
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
