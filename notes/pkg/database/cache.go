package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"notes-manager/configs"
	"notes-manager/envs"
	"notes-manager/models"
	"notes-manager/pkg/enum"
	"notes-manager/pkg/logs"
	"time"

	"github.com/go-redis/redis"
)

type RedisDB struct {
	Cache  *redis.Client
	logger *logs.Logger
}

func NewRedisDB(conf *configs.Config, logs *logs.Logger) (*RedisDB, error) {
	redisURI := fmt.Sprintf("%s:%s", envs.ServerEnvs.RedisHost, envs.ServerEnvs.RedisPort)

	logs.WriteToLog(fmt.Sprint("Redis URI:", redisURI), enum.InfoMsg)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})

	if echoResp := RedisClient.Ping(); echoResp.Val() != "PONG" {
		err := errors.New(echoResp.Val())
		return nil, err
	}

	logs.WriteToLog("Successful connection to the Radis chache", enum.InfoMsg)

	return &RedisDB{
		Cache:  RedisClient,
		logger: logs,
	}, nil
}

func (cache *RedisDB) RecordToChache(notes []models.Note, authorId int) error {
	notesJson, err := json.Marshal(notes)
	if err != nil {
		cache.logger.WriteToLog(fmt.Sprint("Unable to marshaling notes list for caching:", err.Error()), enum.ErrorMsg)
		return err
	}

	collection := fmt.Sprintf("notes/%d", authorId)
	if err := cache.Cache.Set(collection, string(notesJson), time.Minute*1440).Err(); err != nil {
		cache.logger.WriteToLog(fmt.Sprint("Unable to caching notes list", err.Error()), enum.ErrorMsg)
		return err
	}

	cache.logger.WriteToLog("Notes list is cached", enum.InfoMsg)
	return nil
}

func (cache *RedisDB) ReadFromCache(authorId uint) ([]models.Note, error) {
	collection := fmt.Sprintf("notes/%d", authorId)

	records, err := cache.Cache.Get(collection).Result()
	if err == redis.Nil {
		if cache.logLevel == configs.DebugLog {
			log.Println("Cache is empty:", err.Error())
		}
		return nil, err
	} else if err != nil {

	} else if cache.logLevel == configs.DebugLog {
		log.Println("Reading data from cache")
	}

	var notes []models.Note
	err = json.Unmarshal([]byte(records), &notes)
	if err != nil {
		log.Fatalln("Unable unmarshaling records from cache:", err.Error())
		return nil, err
	}

	return notes, nil
}

func (cache *RedisDB) ResetCache(authorId int) (int64, error) {
	collection := fmt.Sprintf("notes/%d", authorId)

	result, err := cache.Cache.Del(collection).Result()
	if err != nil {
		log.Fatalf("Cache reset error:", err.Error())
		return result, err
	} else if cache.logLevel == configs.DebugLog {
		log.Printf("Cache has been reset. %d records deleted", result)
	}
	return result, nil
}

func (cache *RedisDB) writeToLog(msg)

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
