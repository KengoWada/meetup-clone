// Package config provides application configuration settings
// that are used throughout the application. It centralizes the management
// of environment-specific values and parameters for consistency and ease of maintenance.
package config

import (
	"slices"
	"sync"

	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/rs/zerolog"
)

var (
	once      sync.Once
	appConfig Config
)

// Get creates a new Config instance with default values.
func Get() Config {
	once.Do(func() {
		environment := AppEnv(utils.EnvGetString("SERVER_ENVIRONMENT", string(AppEnvProd)))
		if !slices.Contains(Environments, environment) {
			environment = AppEnvProd
		}

		dbAddr := "DB_ADDR"
		if environment == AppEnvTest {
			dbAddr = "TEST_DB_ADDR"
		}

		loglevel := utils.EnvGetInt("LOG_LEVEL", int(zerolog.InfoLevel))
		if environment == AppEnvTest {
			loglevel = int(zerolog.Disabled)
		}

		appConfig = Config{
			Addr:        utils.EnvGetString("SERVER_ADDR", ""),
			Debug:       utils.EnvGetBool("DEBUG", false),
			Environment: environment,
			FrontendURL: utils.EnvGetString("FRONTEND_URL", ""),
			ApiURL:      utils.EnvGetString("API_URL", ""),
			LogLevel:    loglevel,
			SecretKey:   utils.EnvGetString("SECRET_KEY", ""),
			DBConfig: DBConfig{
				Addr:         utils.EnvGetString(dbAddr, ""),
				MaxOpenConns: utils.EnvGetInt("DB_MAX_OPEN_CONNS", 30),
				MaxIdleConns: utils.EnvGetInt("DB_MAX_IDLE_CONNS", 30),
				MaxIdleTime:  utils.EnvGetString("DB_MAX_IDLE_TIME", "15m"),
			},
			AuthConfig: AuthConfig{
				Secret:   utils.EnvGetString("JWT_SECRET_KEY", ""),
				Issuer:   utils.EnvGetString("JWT_ISSUER", "meetup_clone"),
				Audience: utils.EnvGetString("JWT_AUDIENCE", "meetup_clone"),
				Exp:      utils.EnvGetInt("JWT_ACCESS_EXP", 3),
			},
			CacheConfig: CacheConfig{
				ConnURLs: utils.EnvGetStringSlice("MEMCACHED_CONNS", []string{"localhost:11211"}),
				Enabled:  environment != AppEnvTest,
			},
		}
	})

	return appConfig
}
