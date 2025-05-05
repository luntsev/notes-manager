package envs

import (
	"os"

	"github.com/joho/godotenv"
)

type Envs struct {
	MongoHost string
	MongoPort string
	MongoUser string
	MongoPass string
	NotesPort string

	RedisHost string
	RedisPort string
}

var ServerEnvs Envs

func LoadEnvs() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	ServerEnvs.MongoHost = os.Getenv("MONGO_HOST")
	ServerEnvs.MongoPort = os.Getenv("MONGO_PORT")
	ServerEnvs.MongoUser = os.Getenv("MONGO_USER")
	ServerEnvs.MongoPass = os.Getenv("MONGO_PASSWORD")
	ServerEnvs.NotesPort = os.Getenv("NOTES_PORT")

	ServerEnvs.RedisHost = os.Getenv("REDIS_HOST")
	ServerEnvs.RedisPort = os.Getenv("REDIS_PORT")

	return nil
}
