package configs

import (
	"os"
)

type Config struct {
	LogLevel         int
	PostgresDataBase DbConfig
	Jwt              JwtConfig
}

type DbConfig struct {
	PostgresHost     string
	PostgresDbName   string
	PostgresUser     string
	PostgresPassword string
	PostgresPort     string
	PostgresUseSSL   string
}

type JwtConfig struct {
	JwtSecret string
}

func LoadConfig() (*Config, error) {
	conf := &Config{
		PostgresDataBase: DbConfig{
			PostgresHost:     os.Getenv("POSTGRES_HOST"),
			PostgresDbName:   os.Getenv("POSTGRES_DB_NAME"),
			PostgresUser:     os.Getenv("POSTGRES_USER"),
			PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
			PostgresPort:     os.Getenv("POSTGRES_PORT"),
			PostgresUseSSL:   os.Getenv("POSTGRES_USE_SSL"),
		},
		Jwt: JwtConfig{
			JwtSecret: os.Getenv("JWT_SECRET"),
		},
	}

	logLevelStr := os.Getenv("LOG_LEVEL")
	switch logLevelStr {
	case "debug":
		conf.LogLevel = enum.Debug
	case "info":
		conf.LogLevel = enum.Normal
	case "silent":
		conf.LogLevel = enum.Silent
	default:
		return nil, errors.New("logging level is set incorrectly")
	}

}
