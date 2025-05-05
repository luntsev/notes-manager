package repository

import (
	"context"
	"errors"
	"fmt"
	"notes-manager/configs"
	"notes-manager/models"
	"notes-manager/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
)

type NotesRepository struct {
	Db     *database.MondoDB
	Config *configs.Config
}

func NewNotesRepository(config *configs.Config) *NotesRepository {
	return &NotesRepository{
		Config: config,
		Db:     database.NewMongoDB(config),
	}
}

func (repo *NotesRepository) Create(note *models.Note, authorId int, ctx context.Context) (any, error) {
	collection := repo.Db.Client.Database(repo.Config.Db.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	result, err := collection.InsertOne(ctx, note)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (repo *NotesRepository) Read(id string, authorId int, ctx context.Context) (*models.Note, error) {
	collection := repo.Db.Client.Database(repo.Config.Db.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	var note models.Note
	filter := bson.M{"id": id}

	if err := collection.FindOne(ctx, filter).Decode(&note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (repo *NotesRepository) Update(id string, authorId int, note *models.Note, ctx context.Context) (any, error) {
	collection := repo.Db.Client.Database(repo.Config.Db.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	updateFilds := bson.M{}
	if note.Name != nil {
		updateFilds["name"] = note.Name
	}
	if note.Content != nil {
		updateFilds["content"] = note.Content
	}
	update := bson.M{"$set": updateFilds}

	filter := bson.M{"id": id}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, errors.New("notes not found")
	}

	return result, nil
}

func (repo *NotesRepository) Delete(id string, authorId int, ctx context.Context) (any, error) {
	collection := repo.Db.Client.Database(repo.Config.Db.MongoDbName).Collection(fmt.Sprintf("notes/%d", authorId))

	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("notes not found")
	}

	return result, nil
}

func (repo *NotesRepository) List(authorId int, ctx context.Context) (*[]models.Note, error) {
	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", authorId))

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
