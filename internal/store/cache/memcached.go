package cache

import "github.com/bradfitz/gomemcache/memcache"

func NewMemcachedClient(connAddr []string) (*memcache.Client, error) {
	cache := memcache.New(connAddr...)
	if err := cache.Ping(); err != nil {
		return nil, err
	}

	return cache, nil
}
