package config

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
	CacheConfig CacheConfig
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

type CacheConfig struct {
	Enabled  bool
	ConnURLs []string
}
