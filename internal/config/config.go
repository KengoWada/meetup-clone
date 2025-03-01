package config

import (
	"slices"

	"github.com/KengoWada/meetup-clone/internal/utils"
	"github.com/rs/zerolog"
)

const (
	AppEnvDev  AppEnv = "dev"
	AppEnvTest AppEnv = "test"
	AppEnvProd AppEnv = "prod"
)

var Environments = []AppEnv{AppEnvDev, AppEnvTest, AppEnvProd}

type AppEnv string

type Config struct {
	Addr        string
	Debug       bool
	Environment AppEnv
	FrontendURL string
	ApiURL      string
	LogLevel    int
	SecretKey   string
	DBConfig    DBConfig
	AuthConfig  AuthConfig
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type AuthConfig struct {
	Secret   string
	Issuer   string
	Audience string
	Exp      int
}

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
