package database

import (
	"context"
	"fmt"
	"log"
	"notes-manager/configs"
	"notes-manager/envs"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client

type MondoDB struct {
	Client *mongo.Client
}

func NewMongoDB(config *configs.Config) *MondoDB {
	mongoURI := fmt.Sprintf("mongo://%s:%s@%s:%s",
		config.Db.MongoUser,
		config.Db.MongoPass,
		config.Db.MongoHost,
		config.Db.MongoPort)

	if config.LogLevel == configs.DebugLog {
		log.Println("MongoDB URI:", mongoURI)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	err = mongo.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	return &MondoDB{
		Client: mongo,
	}
}

func InitDatabase() error {
	env := &envs.ServerEnvs
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", env.MongoUser, env.MongoPass, env.MongoHost, env.MongoPort)
	log.Println("URI: " + mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI)); err != nil {
		return err
	} else {
		MongoClient = mongo
	}

	err := MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}
