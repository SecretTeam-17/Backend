package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName  string = "sitters"
	colName string = "gamesessions"
)

// Storage - пул подключений к БД.
type Storage struct {
	db *mongo.Client
}

// New - конструктор пула подключений к БД.
func New(path string) (*Storage, error) {
	const operation = "storage.mongodb.New"

	opts := options.Client().ApplyURI(path)
	db, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	collection := db.Database(dbName).Collection(colName)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

// Close - обертка для закрытия подключения к БД.
func (s *Storage) Close() error {
	return s.db.Disconnect(context.Background())
}
