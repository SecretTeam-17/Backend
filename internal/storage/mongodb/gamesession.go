package mongodb

import (
	"context"
	"fmt"
	"petsittersGameServer/internal/storage"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Количество модулей игры зашито в константу пока не ясно, как получать это число динамически.
// const modulesCount int = 4

// CreateSession - создает в базе данных нового юзера и игровую сессию для него.
func (s *Storage) CreateSession(ctx context.Context, name, email string) (*storage.GameSession, error) {
	const operation = "storage.mongodb.CreateSession"

	bsn := bson.D{
		{Key: "_id", Value: primitive.NewObjectID()},
		{Key: "username", Value: name},
		{Key: "email", Value: email},
		{Key: "createdAt", Value: time.Now().UTC()},
		{Key: "updatedAt", Value: time.Now().UTC()},
		{Key: "stats", Value: bson.D{}},
		{Key: "modules", Value: bson.A{}},
		{Key: "minigames", Value: bson.A{}},
	}

	collection := s.db.Database(dbName).Collection(colName)
	resId, err := collection.InsertOne(ctx, bsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	id := resId.InsertedID

	var gs storage.GameSession
	res := collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	err = res.Decode(&gs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// GetSessionById - возвращает игровую сессию из БД по ее Id.
func (s *Storage) GetSessionById(ctx context.Context, id primitive.ObjectID) (*storage.GameSession, error) {
	const operation = "storage.mongodb.GetSessionById"
	var gs storage.GameSession

	// Получаем ссылку на коллекцию, создаем фильтр и ищем игровую сессию в БД.
	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "_id", Value: id}}
	err := collection.FindOne(ctx, filter).Decode(&gs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return &gs, nil
}

// GetSessionByEmail - возвращает игровую сессию из БД по емэйлу ее игрока.
func (s *Storage) GetSessionByEmail(ctx context.Context, email string) (*storage.GameSession, error) {
	const operation = "storage.mongodb.GetSessionByEmail"
	var gs storage.GameSession

	// Получаем ссылку на коллекцию, создаем фильтр и ищем игровую сессию в БД.
	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "email", Value: email}}
	err := collection.FindOne(ctx, filter).Decode(&gs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return &gs, nil
}

// GetSessions - возвращает все игровые сессии из БД.
func (s *Storage) GetSessions(ctx context.Context) ([]storage.GameSession, error) {
	const operation = "storage.mongodb.GetSessions"
	var arr []storage.GameSession

	// Получаем ссылку на коллекцию, задаем порядок сортировки и получаем все  игровые сессии.
	collection := s.db.Database(dbName).Collection(colName)
	opts := options.Find().SetSort(bson.D{{Key: "email", Value: 1}}) // Сортируем по возрастанию поля email
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	// Проверяем, есть ли данные в коллекции.
	if !cursor.Next(ctx) {
		cursor.Close(ctx)
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionsEmpty)
	}
	err = cursor.All(ctx, &arr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	fmt.Println(arr)

	return arr, nil
}

// DeleteSessionById - удаляет данные об игроке и игровой сессии по ее id.
func (s *Storage) DeleteSessionById(ctx context.Context, id primitive.ObjectID) error {
	const operation = "storage.mongodb.DeleteSessionById"

	// Получаем ссылку на коллекцию, создаем фильтр и удаляем игровую сессию в БД.
	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "_id", Value: id}}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	return nil
}

// TruncateTables - удаляет все записи из таблиц users и game_sessions.
func (s *Storage) TruncateData(ctx context.Context) error {
	const operation = "storage.mongodb.TruncateData"

	// Получаем ссылку на коллекцию и удаляем все игровые сессии в БД.
	collection := s.db.Database(dbName).Collection(colName)
	res, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", operation, storage.ErrSessionsEmpty)
	}

	return nil
}
