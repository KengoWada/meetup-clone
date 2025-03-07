package cache

import (
	"github.com/KengoWada/meetup-clone/internal/models"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	CacheKeyUser string = "user:"
	CacheTTLUser int32  = 60 * 60 // 1 hour in seconds
)

type CacheKey string

type Store struct {
	Users interface {
		Get(ID int64) (*models.User, error)
		Set(user *models.User) error
		Delete(ID int64) error
	}
}

func NewCacheStore(memcached *memcache.Client) Store {
	return Store{
		Users: &UserStore{cacheDB: memcached},
	}
}
