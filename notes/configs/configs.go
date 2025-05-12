package configs

import (
	"errors"
	"os"
	"strconv"

	"github.com/luntsev/notes-manager/notes/pkg/enum"
)

type Config struct {
	LogLevel      int
	MongoDataBase DbConfig
	RedisCache    CacheConfig
	Jwt           JwtConfig
	Port          int
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

type JwtConfig struct {
	JwtSecret          string
	AccessTokerExpire  int
	RefreshTokenExpire int
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
		Jwt: JwtConfig{
			JwtSecret: os.Getenv("JWT_SECRET"),
		},
	}

	accessTokenExpire, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRE"))
	if err != nil {
		return nil, err
	} else if accessTokenExpire <= 0 {
		return nil, errors.New("bad ACCESS_TOKEN_EXPIRE in envarenment variable")
	}

	refreshTokenExpire, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRE"))
	if err != nil {
		return nil, err
	} else if refreshTokenExpire <= 0 {
		return nil, errors.New("bad REFRESH_TOKEN_EXPIRE in envarenment variable")
	}

	conf.Jwt.AccessTokerExpire = accessTokenExpire
	conf.Jwt.RefreshTokenExpire = refreshTokenExpire

	port, err := strconv.Atoi(os.Getenv("NOTES_PORT"))
	if err != nil {
		return nil, err
	} else if port <= 0 {
		return nil, errors.New("bad port in envarenment variable")
	}
	conf.Port = port

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
	case conf.Jwt.JwtSecret == "":
		return nil, errors.New("the SECRET environment variable is not set")
	default:
		return conf, nil
	}
}
