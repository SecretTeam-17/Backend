package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"petsittersGameServer/internal/storage"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type gsBase struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
}

type gsExtend struct {
	Stats     json.RawMessage `json:"stats" bson:"stats"`
	Modules   json.RawMessage `json:"modules" bson:"modules"`
	Minigames json.RawMessage `json:"minigames" bson:"minigames"`
}

// CreateSession - создает в базе данных нового юзера и игровую сессию для него.
func (s *Storage) CreateSession(ctx context.Context, name, email string, stats, modules, minigames []byte) (*storage.GameSession, error) {
	const operation = "storage.mongodb.CreateSession"

	// Преобразуем слайсы байт в JSON для корректной записи в БД.
	var inputStats, inputMod, inputMini interface{}
	err := bson.UnmarshalExtJSON(stats, true, &inputStats)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = bson.UnmarshalExtJSON(modules, true, &inputMod)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = bson.UnmarshalExtJSON(minigames, true, &inputMini)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	bsn := bson.D{
		{Key: "_id", Value: primitive.NewObjectID()},
		{Key: "username", Value: name},
		{Key: "email", Value: email},
		{Key: "createdAt", Value: time.Now().UTC()},
		{Key: "updatedAt", Value: time.Now().UTC()},
		{Key: "stats", Value: inputStats},
		{Key: "modules", Value: inputMod},
		{Key: "minigames", Value: inputMini},
	}

	collection := s.db.Database(dbName).Collection(colName)
	resId, err := collection.InsertOne(ctx, bsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	id := resId.InsertedID

	res := collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}})

	// rba := bson.NewRegistry()
	// rba.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, reflect.TypeOf(bson.M{}))

	var gsbase gsBase
	err = res.Decode(&gsbase)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	var bs bson.D
	err = res.Decode(&bs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	bytes, err := bson.MarshalExtJSON(bs, true, true)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	var gsext gsExtend
	err = json.Unmarshal(bytes, &gsext)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	fmt.Println("1 ", gsbase)
	fmt.Println("2 ", bs)
	fmt.Println("3 ", string(bytes))
	// fmt.Println("4 ", gsext)

	gs := storage.GameSession{
		Id:        gsbase.Id,
		Username:  gsbase.Username,
		Email:     gsbase.Email,
		CreatedAt: gsbase.CreatedAt,
		UpdatedAt: gsbase.UpdatedAt,
		Stats:     gsext.Stats,
		Modules:   gsext.Modules,
		Minigames: gsext.Minigames,
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
