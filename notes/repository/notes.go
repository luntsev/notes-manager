package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"notes-manager/configs"
	"notes-manager/models"
	"notes-manager/pkg/database"
	"notes-manager/pkg/logs"
)

type NotesRepository struct {
	db     *database.MongoDB
	cache  *database.RedisDB
	conf   *configs.Config
	logger *logs.Logger
}

func NewNotesRepository(config *configs.Config, mongoDB *database.MongoDB, redisDB *database.RedisDB, logs *logs.Logger) *NotesRepository {
	return &NotesRepository{
		db:     mongoDB,
		cache:  redisDB,
		conf:   config,
		logger: logs,
	}
}

func (repo *NotesRepository) Create(note *models.Note, ctx context.Context) (any, error) {
	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", note.AuthorId))

	result, err := collection.InsertOne(ctx, note)
	if err != nil {
		repo.logger.WriteError(fmt.Sprintf("Unable write to database when creating record: %s", err.Error()))
		return nil, err
	}

	go repo.cache.ResetCache(note.AuthorId)

	return result, nil
}

func (repo *NotesRepository) Read(id string, authorId uint, ctx context.Context) (*models.Note, error) {
	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	var note models.Note
	filter := bson.M{"id": id}

	if err := collection.FindOne(ctx, filter).Decode(&note); err != nil {
		repo.logger.WriteError(fmt.Sprintf("Unable read from database: %s", err.Error()))
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
		repo.logger.WriteError(fmt.Sprintf("Unable write to database when updating record: %s", err.Error()))
		return err
	}
	if result.MatchedCount == 0 {
		repo.logger.WriteInfo("Notes not found when reading record")
		return errors.New("notes not found")
	}

	go repo.cache.ResetCache(note.AuthorId)

	return nil
}

func (repo *NotesRepository) Delete(id string, authorId uint, ctx context.Context) (int64, error) {
	collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		repo.logger.WriteError(fmt.Sprintf("Unable deleting record: %s", err.Error()))
		return 0, err
	}

	if result.DeletedCount == 0 {
		repo.logger.WriteWarn("Notes not found when deleting record")
		return 0, errors.New("notes not found")
	}

	go repo.cache.ResetCache(authorId)

	return result.DeletedCount, nil
}

func (repo *NotesRepository) List(authorId uint, ctx context.Context) (*[]models.Note, error) {
	notes, err := repo.cache.ReadFromCache(authorId)
	if err != nil {
		collection := repo.db.Db.Database(repo.conf.MongoDataBase.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			repo.logger.WriteError(fmt.Sprintf("Unable reading records: %s", err.Error()))
			return nil, err
		}
		defer cursor.Close(ctx)

		var readingNotes []models.Note

		if err = cursor.All(ctx, &readingNotes); err != nil {
			repo.logger.WriteError(fmt.Sprintf("Unable reading records: %s", err.Error()))
			return nil, err
		}

		if len(readingNotes) == 0 {
			repo.logger.WriteInfo("Collections is empty")
			return nil, errors.New("collection is empty")
		} else {
			notes = &readingNotes
			go repo.cache.WriteToChache(notes, int(authorId))
		}
	}

	return notes, nil
}
