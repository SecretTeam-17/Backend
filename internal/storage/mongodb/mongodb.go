package mongodb

import (
	"context"
	"fmt"
	"log"
	"petsittersGameServer/internal/config"

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

// New - обертка для конструктора пула подключений new.
func New(cfg *config.Config) *Storage {
	storage, err := new(cfg.StoragePath, cfg.StorageUser, cfg.StoragePasswd)
	if err != nil {
		log.Fatalf("failed to init storage: %s", err.Error())
	}
	return storage
}

// new - конструктор пула подключений к БД.
func new(path, user, password string) (*Storage, error) {
	const operation = "storage.mongodb.New"

	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      user,
		Password:      password,
	}

	// Задаем опции подключения.
	opts := options.Client().ApplyURI(path).SetAuth(credential)
	// Создаем подключение к MongoDB и проверяем его.
	db, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Создаем уникальный индекс по полю email.
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
