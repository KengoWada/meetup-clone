package app

type Config struct {
	Addr        string
	Debug       bool
	Environment string
	FrontendURL string
	DBConfig    DBConfig
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
