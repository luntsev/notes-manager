package configs

import (
	"errors"
	"notes-manager/pkg/enum"
	"os"
)

type Config struct {
	LogLevel      int
	MongoDataBase DbConfig
	RedisCache    CacheConfig
	Auth          AuthConfig
}

type DbConfig struct {
	MongoHost   string
	MongoPort   string
	MongoDbName string
	MongoUser   string
	MongoPass   string
	NotesPort   string
}

type CacheConfig struct {
	RedisHost string
	RedisPort string
}

type AuthConfig struct {
	JwtSecret string
}

func LoadConfig() (*Config, error) {
	conf := &Config{
		MongoDataBase: DbConfig{
			MongoHost:   os.Getenv("MONGO_HOST"),
			MongoPort:   os.Getenv("MONGO_PORT"),
			MongoDbName: os.Getenv("MONGO_DB_NAME"),
			MongoUser:   os.Getenv("MONGO_USER"),
			MongoPass:   os.Getenv("MONGO_PASSWORD"),
			NotesPort:   os.Getenv("NOTES_PORT"),
		},
		RedisCache: CacheConfig{
			RedisHost: os.Getenv("REDIS_HOST"),
			RedisPort: os.Getenv("REDIS_PORT"),
		},
		Auth: AuthConfig{
			JwtSecret: os.Getenv("SECRET"),
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

	switch {
	case conf.MongoDataBase.MongoHost == "":
		return nil, errors.New("the MONGO_HOST environment variable is not set")
	case conf.MongoDataBase.MongoPort == "":
		return nil, errors.New("the MONGO_PORT environment variable is not set")
	case conf.MongoDataBase.MongoDbName == "":
		return nil, errors.New("the MONGO_DB_NAME environment variable is not set")
	case conf.MongoDataBase.MongoUser == "":
		return nil, errors.New("the MONGO_USER environment variable is not set")
	case conf.MongoDataBase.MongoPass == "":
		return nil, errors.New("the MONGO_PASSWORD environment variable is not set")
	case conf.MongoDataBase.MongoPort == "":
		return nil, errors.New("the MONGO_PORT environment variable is not set")
	case conf.RedisCache.RedisHost == "":
		return nil, errors.New("the REDIS_HOST environment variable is not set")
	case conf.RedisCache.RedisPort == "":
		return nil, errors.New("the REDIS_PORT environment variable is not set")
	case conf.Auth.JwtSecret == "":
		return nil, errors.New("the SECRET environment variable is not set")
	default:
		return conf, nil
	}
}
