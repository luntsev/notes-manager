package database

import (
	"context"
	"fmt"
	"log"
	"notes-manager/envs"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client

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
