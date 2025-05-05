package database

import (
	"errors"
	"fmt"
	"notes-manager/envs"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis() error {
	redisURI := fmt.Sprintf("%s:%s", envs.ServerEnvs.RedisHost, envs.ServerEnvs.RedisPort)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})

	if echoResp := RedisClient.Ping(); echoResp.Val() != "PONG" {
		err := errors.New(echoResp.Val())
		return err
	}
	return nil
}
