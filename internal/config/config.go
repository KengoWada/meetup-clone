// Package config provides application configuration settings
// that are used throughout the application. It centralizes the management
// of environment-specific values and parameters for consistency and ease of maintenance.
package config

import (
	"slices"

	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/rs/zerolog"
)

// Valid environment values for AppEnv.
const (
	AppEnvDev  AppEnv = "dev"  // Development environment
	AppEnvTest AppEnv = "test" // Test environment
	AppEnvProd AppEnv = "prod" // Production environment
)

// Valid environment values for AppEnv.
var Environments = []AppEnv{AppEnvDev, AppEnvTest, AppEnvProd}

// AppEnv represents the application environment as a string.
// It ensures the environment is one of the following valid values:
//   - "dev"  → Development environment
//   - "test" → Testing environment
//   - "prod" → Production environment
type AppEnv string

// Config holds the application configuration settings.
type Config struct {
	Addr        string     // The application port in the format ":8000".
	Debug       bool       // Run the application in debug mode if true.
	Environment AppEnv     // The environment the app is running in (e.g., "dev", "test", "prod").
	FrontendURL string     // The URL of the frontend application.
	ApiURL      string     // The API URL, should match the Addr. (e.g., "localhost:8000").
	LogLevel    int        // The log level for the app.
	SecretKey   string     // The secret key used for generating and signing tokens.
	DBConfig    DBConfig   // The application database configurations
	AuthConfig  AuthConfig // The application authentication configurations.
}

// DBConfig holds the database connection configuration settings.
type DBConfig struct {
	Addr         string // The database connection string.
	MaxOpenConns int    // The maximum number of open connections to the database.
	MaxIdleConns int    // The maximum number of idle connections to keep in the pool.
	MaxIdleTime  string // The maximum amount of time a connection can remain idle (e.g., "5m", "30s").
}

// AuthConfig holds the configuration settings for authentication and JWT handling.
type AuthConfig struct {
	Secret   string // The secret key used to sign and verify JWT tokens.
	Issuer   string // The issuer claim (iss) for the tokens.
	Audience string // The audience claim (aud) for the tokens.
	Exp      int    // The token expiration time in hours.
}

// NewConfig creates a new Config instance with default values.
//
// Returns:
//   - Config: a Config struct with initialized default settings.
func NewConfig() Config {
	environment := AppEnv(utils.GetString("SERVER_ENVIRONMENT", string(AppEnvProd)))
	if !slices.Contains(Environments, environment) {
		environment = AppEnvProd
	}

	dbAddr := "DB_ADDR"
	if environment == AppEnvTest {
		dbAddr = "TEST_DB_ADDR"
	}

	cfg := Config{
		Addr:        utils.GetString("SERVER_ADDR", ""),
		Debug:       utils.GetBool("DEBUG", false),
		Environment: environment,
		FrontendURL: utils.GetString("FRONTEND_URL", ""),
		ApiURL:      utils.GetString("API_URL", ""),
		LogLevel:    utils.GetInt("LOG_LEVEL", int(zerolog.InfoLevel)),
		SecretKey:   utils.GetString("SECRET_KEY", ""),
		DBConfig: DBConfig{
			Addr:         utils.GetString(dbAddr, ""),
			MaxOpenConns: utils.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: utils.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  utils.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		AuthConfig: AuthConfig{
			Secret:   utils.GetString("JWT_SECRET_KEY", ""),
			Issuer:   utils.GetString("JWT_ISSUER", "meetup_clone"),
			Audience: utils.GetString("JWT_AUDIENCE", "meetup_clone"),
			Exp:      utils.GetInt("JWT_ACCESS_EXP", 3),
		},
	}

	return cfg
}
