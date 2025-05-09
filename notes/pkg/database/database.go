package database

import (
	"context"
	"fmt"
	"github.com/luntsev/notes-manager/notes/configs"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoDB struct {
	Db     *mongo.Client
	Logger *logs.Logger
}

func NewMongoDB(conf *configs.Config, logger *logs.Logger) *MongoDB {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		conf.MongoDataBase.MongoUser,
		conf.MongoDataBase.MongoPass,
		conf.MongoDataBase.MongoHost,
		conf.MongoDataBase.MongoPort)

	logger.WriteInfo(fmt.Sprintf("MongoDB URI: %s", mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.WriteError(fmt.Sprintf("Unable connect to database: %s", err.Error()))
		panic(err)
	}

	err = mongo.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.WriteError(fmt.Sprintf("Unable connect to database: %s", err.Error()))
		panic(err)
	}

	logger.WriteInfo("Successful connection to the database")

	return &MongoDB{
		Db:     mongo,
		Logger: logger,
	}
}
