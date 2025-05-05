package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	DebugLog = "debug"
)

type Config struct {
	LogLevel string
	Db       DbConfig
	Cache    CacheConfig
	Auth     AuthConfig
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

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("File \".env\" not foud: %s. Using default config.", err.Error())
	}

	return &Config{
		LogLevel: os.Getenv("LOG_LEVEL"),
		Db: DbConfig{
			MongoHost:   os.Getenv("MONGO_HOST"),
			MongoPort:   os.Getenv("MONGO_PORT"),
			MongoDbName: os.Getenv("MONGO_DB_NAME"),
			MongoUser:   os.Getenv("MONGO_USER"),
			MongoPass:   os.Getenv("MONGO_PASSWORD"),
			NotesPort:   os.Getenv("NOTES_PORT"),
		},
		Cache: CacheConfig{
			RedisHost: os.Getenv("REDIS_HOST"),
			RedisPort: os.Getenv("REDIS_PORT"),
		},
		Auth: AuthConfig{
			JwtSecret: os.Getenv("SECRET"),
		},
	}
}
