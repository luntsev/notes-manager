package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"notes-manager/configs"
	"notes-manager/envs"
	"notes-manager/models"
	"notes-manager/pkg/logs"
	"time"

	"github.com/go-redis/redis"
)

type RedisDB struct {
	Cache  *redis.Client
	logger *logs.Logger
}

func NewRedisDB(conf *configs.Config, logger *logs.Logger) (*RedisDB, error) {
	redisURI := fmt.Sprintf("%s:%s", envs.ServerEnvs.RedisHost, envs.ServerEnvs.RedisPort)

	logger.WriteInfo(fmt.Sprintf("Redis URI: %s", redisURI))

	RedisClient := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})

	if echoResp := RedisClient.Ping(); echoResp.Val() != "PONG" {
		err := errors.New(echoResp.Val())
		logger.WriteError(fmt.Sprintf("Unable connect to Radis cache: %s", err.Error()))
		return nil, err
	}

	logger.WriteInfo("Successful connection to the Radis chache")

	return &RedisDB{
		Cache:  RedisClient,
		logger: logger,
	}, nil
}

func (cache *RedisDB) WriteToChache(notes *[]models.Note, authorId int) error {
	notesJson, err := json.Marshal(notes)
	if err != nil {
		cache.logger.WriteError(fmt.Sprintf("Unable to marshaling notes list for caching: %s", err.Error()))
		return err
	}

	collection := fmt.Sprintf("notes/%d", authorId)
	if err := cache.Cache.Set(collection, string(notesJson), time.Minute*1440).Err(); err != nil {
		cache.logger.WriteError(fmt.Sprintf("Unable to caching notes list: %s", err.Error()))
		return err
	}

	cache.logger.WriteInfo("Notes list is cached")
	return nil
}

func (cache *RedisDB) ReadFromCache(authorId uint) (*[]models.Note, error) {
	collection := fmt.Sprintf("notes/%d", authorId)

	records, err := cache.Cache.Get(collection).Result()
	if err == redis.Nil {
		cache.logger.WriteWarn(fmt.Sprintf("Cache is empty: %s", err.Error()))
		return nil, err
	} else if err != nil {
		cache.logger.WriteError(fmt.Sprintf("Cache readint error: %s", err.Error()))
		return nil, err
	}
	cache.logger.WriteInfo("Reading data from cache")

	var notes []models.Note
	err = json.Unmarshal([]byte(records), &notes)
	if err != nil {
		cache.logger.WriteError(fmt.Sprintf("Unable unmarshaling records from cache:", err.Error()))
		return nil, err
	}

	return &notes, nil
}

func (cache *RedisDB) ResetCache(authorId uint) (int64, error) {
	collection := fmt.Sprintf("notes/%d", authorId)

	result, err := cache.Cache.Del(collection).Result()
	if err != nil {
		cache.logger.WriteError(fmt.Sprintf("Cache reset error: %s", err.Error()))
		return result, err
	}
	cache.logger.WriteInfo(fmt.Sprintf("Cache has been reset. %d records deleted", result))
	return result, nil
}
