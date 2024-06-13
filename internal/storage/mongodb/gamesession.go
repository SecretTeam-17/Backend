package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"petsittersGameServer/internal/storage"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func buildGS(res *mongo.SingleResult) (*storage.GameSession, error) {
	const operation = "storage.mongodb.buildGS"

	var gsbase gsBase
	err := res.Decode(&gsbase)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	var bs bson.D
	err = res.Decode(&bs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	bytes, err := bson.MarshalExtJSON(bs, false, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	var gsext gsExtend
	err = json.Unmarshal(bytes, &gsext)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

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

// CreateSession - создает в базе данных нового юзера и игровую сессию для него.
func (s *Storage) CreateSession(ctx context.Context, name, email string, stats, modules, minigames []byte) (*storage.GameSession, error) {
	const operation = "storage.mongodb.CreateSession"

	// Преобразуем слайсы байт в JSON для корректной записи в БД.
	var inputStats, inputMod, inputMini interface{}
	err := bson.UnmarshalExtJSON(stats, false, &inputStats)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = bson.UnmarshalExtJSON(modules, false, &inputMod)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = bson.UnmarshalExtJSON(minigames, false, &inputMini)
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

	gs, err := buildGS(res)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return gs, nil
}

// GetSessionById - возвращает игровую сессию из БД по ее Id.
func (s *Storage) GetSessionById(ctx context.Context, id primitive.ObjectID) (*storage.GameSession, error) {
	const operation = "storage.mongodb.GetSessionById"

	// Получаем ссылку на коллекцию, создаем фильтр и ищем игровую сессию в БД.
	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "_id", Value: id}}
	res := collection.FindOne(ctx, filter)

	gs, err := buildGS(res)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return gs, nil
}

// GetSessionByEmail - возвращает игровую сессию из БД по емэйлу ее игрока.
func (s *Storage) GetSessionByEmail(ctx context.Context, email string) (*storage.GameSession, error) {
	const operation = "storage.mongodb.GetSessionByEmail"

	// Получаем ссылку на коллекцию, создаем фильтр и ищем игровую сессию в БД.
	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "email", Value: email}}
	res := collection.FindOne(ctx, filter)

	gs, err := buildGS(res)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return gs, nil
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
	defer cursor.Close(ctx)

	// Аналог метода buildGS для курсора
	for cursor.Next(ctx) {

		var gsbase gsBase
		err := cursor.Decode(&gsbase)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}

		var bs bson.D
		err = cursor.Decode(&bs)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
		bytes, err := bson.MarshalExtJSON(bs, false, false)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
		var gsext gsExtend
		err = json.Unmarshal(bytes, &gsext)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}

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
		arr = append(arr, gs)
	}

	// Проверяем, есть ли данные в коллекции.
	if len(arr) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionsEmpty)
	}

	return arr, nil
}

// CleanSession - устанавливает изменяемые поля игровой сессии в значение по-умолчанию.
func (s *Storage) CleanSession(ctx context.Context, id primitive.ObjectID) (*storage.GameSession, error) {
	const operation = "storage.mongodb.CleanSession"

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "updatedAt", Value: time.Now().UTC()},
			{Key: "stats", Value: nil},
			{Key: "modules", Value: nil},
			{Key: "minigames", Value: nil},
		}},
	}
	filter := bson.D{{Key: "_id", Value: id}}

	collection := s.db.Database(dbName).Collection(colName)
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if res.ModifiedCount == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	found := collection.FindOne(ctx, filter)

	gs, err := buildGS(found)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return gs, nil
}

// UpdateSession - устанавливает изменяемые поля игровой сессии в соответствии с переданными значениями.
func (s *Storage) UpdateSession(ctx context.Context, gs storage.GameSession) error {
	const operation = "storage.mongodb.UpdateSession"

	// Преобразуем слайсы байт в JSON для корректной записи в БД.
	var inputStats, inputMod, inputMini interface{}
	err := bson.UnmarshalExtJSON(gs.Stats, true, &inputStats)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	err = bson.UnmarshalExtJSON(gs.Modules, true, &inputMod)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	err = bson.UnmarshalExtJSON(gs.Minigames, true, &inputMini)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "updatedAt", Value: time.Now().UTC()},
			{Key: "stats", Value: inputStats},
			{Key: "modules", Value: inputMod},
			{Key: "minigames", Value: inputMini},
		}},
	}
	filter := bson.D{{Key: "_id", Value: gs.Id}}

	collection := s.db.Database(dbName).Collection(colName)
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	return nil
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
