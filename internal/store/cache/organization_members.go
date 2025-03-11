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

func (s *OrganizationMemberStore) Get(userID, orgID int64) (*models.OrganizationMember, error) {
	cacheKey := fmt.Sprintf("%s:%d,%d", CacheKeyOrgMember, userID, orgID)

	item, err := getFromCache(s.cacheDB, cacheKey)
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
	cacheKey := fmt.Sprintf("%s:%d,%d", CacheKeyOrgMember, member.UserProfileID, member.OrganizationID)

	orgBytes, err := json.Marshal(member)
	if err != nil {
		return err
	}

	roleItem := &memcache.Item{Key: cacheKey, Value: orgBytes, Expiration: CacheTTLOrg}
	if err := s.cacheDB.Set(roleItem); err != nil {
		return err
	}

	return nil
}

func (s *OrganizationMemberStore) Delete(userID, orgID int64) error {
	cacheKey := fmt.Sprintf("%s:%d,%d", CacheKeyOrgMember, userID, orgID)

	err := s.cacheDB.Delete(cacheKey)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	return nil
}
