package database

import (
	"context"
	"fmt"
	"log"
	"notes-manager/configs"
	"notes-manager/envs"
	"notes-manager/pkg/enum"
	"notes-manager/pkg/logs"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client

type MongoDB struct {
	Db     *mongo.Client
	Logger *logs.Logger
}

func NewMongoDB(conf *configs.Config, logger *logs.Logger) *MongoDB {
	mongoURI := fmt.Sprintf("mongo://%s:%s@%s:%s",
		conf.MongoDataBase.MongoUser,
		conf.MongoDataBase.MongoPass,
		conf.MongoDataBase.MongoHost,
		conf.MongoDataBase.MongoPort)

	logger.WriteToLog(fmt.Sprint("MongoDB URI:", mongoURI), enum.InfoMsg)

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

	logger.WriteToLog("Successful connection to the database", enum.InfoMsg)

	return &MongoDB{
		Db:     mongo,
		Logger: logger,
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
