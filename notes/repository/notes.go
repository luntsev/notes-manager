package repository

import (
	"context"
	"errors"
	"fmt"
	"notes-manager/configs"
	"notes-manager/models"
	"notes-manager/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotesRepository struct {
	db    *database.MongoDB
	cache *database.RedisDB
	conf  *configs.Config
}

func NewNotesRepository(config *configs.Config, mongoDB *database.MongoDB, redisDB *database.RedisDB) *NotesRepository {
	return &NotesRepository{
		db:    mongoDB,
		cache: redisDB,
		conf:  config,
	}
}

func (repo *NotesRepository) Create(note *models.Note, ctx context.Context) (any, error) {
	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", note.AuthorId))

	result, err := collection.InsertOne(ctx, note)
	if err != nil {
		return nil, err
	}

	go repo.cache.ResetCache(note.AuthorId)

	return result, err
}

func (repo *NotesRepository) Read(id string, authorId uint, ctx context.Context) (*models.Note, error) {
	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	var note models.Note
	filter := bson.M{"id": id}

	if err := collection.FindOne(ctx, filter).Decode(&note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (repo *NotesRepository) Update(note *models.Note, ctx context.Context) error {
	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", note.AuthorId))

	updateFilds := bson.M{}
	if note.Name != nil {
		updateFilds["name"] = note.Name
	}
	if note.Content != nil {
		updateFilds["content"] = note.Content
	}
	update := bson.M{"$set": updateFilds}

	filter := bson.M{"id": note.Id}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("notes not found")
	}

	return nil
}

func (repo *NotesRepository) Delete(id string, authorId uint, ctx context.Context) (int64, error) {
	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	if result.DeletedCount == 0 {
		return 0, errors.New("notes not found")
	}

	return result.DeletedCount, nil
}

func (repo *NotesRepository) List(authorId uint, ctx context.Context) (*[]models.Note, error) {
	notes, err := repo.cache.ReadFromCache(authorId)
	if err != nil {

	}

	cache, err := database.RedisClient.Get(fmt.Sprintf("notes/%d", authorId)).Result()

	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var (
		note  models.Note
		notes []models.Note
	)

	for cursor.Next(ctx) {
		err := cursor.Decode(&note)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if len(notes) == 0 {
		return nil, errors.New("collection is empty")
	}

	return &notes, nil
}
