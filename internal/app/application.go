package app

import (
	"database/sql"

	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/db"
	"github.com/KengoWada/meetup-clone/internal/store"
	"github.com/KengoWada/meetup-clone/internal/store/cache"
	"github.com/bradfitz/gomemcache/memcache"
)

type Application struct {
	Config        config.Config
	Store         store.Store
	CacheStore    cache.Store
	Authenticator auth.Authenticator
}

type AppItems struct {
	App       *Application
	DB        *sql.DB
	Memcached *memcache.Client
}

func NewApplication() (*AppItems, error) {
	cfg := config.Get()

	appItems := &AppItems{}

	// Make database connection
	db, err := db.New(
		cfg.DBConfig.Addr,
		cfg.DBConfig.MaxOpenConns,
		cfg.DBConfig.MaxIdleConns,
		cfg.DBConfig.MaxIdleTime,
		cfg.Environment,
	)
	if err != nil {
		return appItems, err
	}
	l.Info().Msg("successfully connected to postgres")
	appItems.DB = db

	// Cache connection
	var memcached *memcache.Client
	if cfg.CacheConfig.Enabled {
		memcached, err = cache.NewMemcachedClient(cfg.CacheConfig.ConnURLs)
		if err != nil {
			return appItems, err
		}
		l.Info().Msg("successfully connected to memcached")
		appItems.Memcached = memcached
	}

	// Create JWT Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.AuthConfig.Secret, cfg.AuthConfig.Audience, cfg.AuthConfig.Issuer)

	// Create Global App Store
	store := store.NewStore(db)
	cacheStore := cache.NewCacheStore(memcached)

	app := &Application{
		Config:        cfg,
		Store:         store,
		CacheStore:    cacheStore,
		Authenticator: jwtAuthenticator,
	}
	appItems.App = app

	return appItems, nil
}
