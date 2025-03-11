package cache

import (
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	CacheKeyUser string = "user"
	CacheTTLUser int32  = 60 * 60 // 1 hour in seconds

	CacheKeyOrg string = "org"
	CacheTTLOrg int32  = 60 * 60 // 1 hour in seconds

	CacheKeyRole string = "role"
	CacheTTLRole int32  = 60 * 60 // 1 hour in seconds

	CacheKeyOrgMember string = "org_member"
	CacheTTLOrgMember int32  = 60 * 60 // 1 hour in seconds
)

type CacheKey string

type Store struct {
	Users interface {
		Get(ID int64) (*models.User, error)
		Set(user *models.User) error
		Delete(ID int64) error
	}
	Organizations interface {
		Get(ID int64) (*models.Organization, error)
		Set(organization *models.Organization) error
		Delete(ID int64) error
	}
	Roles interface {
		Get(ID int64) (*models.Role, error)
		Set(role *models.Role) error
		Delete(ID int64) error
	}
}

func NewCacheStore(memcached *memcache.Client) Store {
	return Store{
		Users:         &UserStore{cacheDB: memcached},
		Organizations: &OrganizationStore{cacheDB: memcached},
		Roles:         &RoleStore{cacheDB: memcached},
	}
}

func getFromCache(cache *memcache.Client, cacheKey string) (*memcache.Item, error) {
	item, err := cache.Get(cacheKey)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}

		return nil, err
	}

	return item, nil
}
