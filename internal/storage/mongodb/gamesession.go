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

// gsBase - структура для декодинга простых данных.
type gsBase struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
}

// gsExtend - структура для декодинга данных в json формате.
type gsExtend struct {
	Stats     json.RawMessage `json:"stats" bson:"stats"`
	Modules   json.RawMessage `json:"modules" bson:"modules"`
	Minigames json.RawMessage `json:"minigames" bson:"minigames"`
}

// buildGS - вспомогательная функция, строит из структур gsBase и gsExtend одну структуру GameSession.
func buildGS(res *mongo.SingleResult) (*storage.GameSession, error) {
	const operation = "storage.mongodb.buildGS"

	// Декодируем простые данные из ответа БД.
	var gsbase gsBase
	err := res.Decode(&gsbase)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Декодируем json данные из ответа БД.
	var bs bson.D
	err = res.Decode(&bs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	bytes, err := bson.MarshalExtJSON(bs, false, false)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Записываем json данные в структуру gsExtend.
	var gsext gsExtend
	err = json.Unmarshal(bytes, &gsext)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Создаем структуру GameSession и заполняем ее данными из двух структур gsBase и gsExtend.
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

// CreateSession - создает в базе данных новую игровую сессию.
func (s *Storage) CreateSession(ctx context.Context, name, email string, stats, modules, minigames []byte) (*storage.GameSession, error) {
	const operation = "storage.mongodb.CreateSession"

	// При POST запросе на создание нового игрока проверяем email. Если такой игрок уже сущетсвует,
	// то возвращаем его игровую сессию.
	gs, err := s.GetSessionByEmail(ctx, email)
	if err == nil {
		return gs, storage.ErrUserExists
	}

	// Преобразуем слайсы байт json.RawMessage для корректной записи в БД.
	var inputStats, inputMod, inputMini interface{}
	err = bson.UnmarshalExtJSON(stats, false, &inputStats)
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

	// Создаем bson.D структуру для записи в БД.
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

	// Получаем ссылку на коллекцию и записываем в нее данные новой игровой сессии.
	collection := s.db.Database(dbName).Collection(colName)
	resId, err := collection.InsertOne(ctx, bsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	id := resId.InsertedID

	// Получаем игровую сессию из БД.
	res := collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	gs, err = buildGS(res)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return gs, nil
}

// GetSessionById - возвращает игровую сессию из БД по ее Id.
func (s *Storage) GetSessionById(ctx context.Context, id string) (*storage.GameSession, error) {
	const operation = "storage.mongodb.GetSessionById"

	// Преобразуем строку в ObjectId.
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrInput)
	}
	// Получаем ссылку на коллекцию, создаем фильтр и ищем игровую сессию в БД.
	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "_id", Value: idObj}}

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

	// Получаем ссылку на коллекцию, задаем порядок сортировки и получаем курсор на все игровые сессии.
	collection := s.db.Database(dbName).Collection(colName)
	opts := options.Find().SetSort(bson.D{{Key: "email", Value: 1}}) // Сортируем по возрастанию поля email.
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	defer cursor.Close(ctx)

	// Аналог метода buildGS для курсора. Перебираем все игровые сессии из БД и записываем их с слайс.
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

	// Проверяем, была ли получена хотя бы одна игровая сессия.
	if len(arr) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionsEmpty)
	}

	return arr, nil
}

// CleanSession - устанавливает изменяемые поля игровой сессии в значение по-умолчанию.
func (s *Storage) CleanSession(ctx context.Context, id string) (*storage.GameSession, error) {
	const operation = "storage.mongodb.CleanSession"

	// Преобразуем строку в ObjectId.
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrInput)
	}

	// Задаем поля для обновления/очистки и их значения.
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "updatedAt", Value: time.Now().UTC()},
			{Key: "stats", Value: nil},
			{Key: "modules", Value: nil},
			{Key: "minigames", Value: nil},
		}},
	}
	// Создаем фильтр и обновляем необходимые данные в БД.
	filter := bson.D{{Key: "_id", Value: idObj}}
	collection := s.db.Database(dbName).Collection(colName)
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if res.ModifiedCount == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	// Получаем игровую сессию из БД.
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

	// Преобразуем слайсы байт json.RawMessage для корректной записи в БД.
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

	// Задаем поля для обновления и их значения.
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "updatedAt", Value: time.Now().UTC()},
			{Key: "stats", Value: inputStats},
			{Key: "modules", Value: inputMod},
			{Key: "minigames", Value: inputMini},
		}},
	}

	// Создаем фильтр и обновляем необходимые данные в БД.
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
func (s *Storage) DeleteSessionById(ctx context.Context, id string) error {
	const operation = "storage.mongodb.DeleteSessionById"

	// Преобразуем строку в ObjectId.
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, storage.ErrInput)
	}

	// Получаем ссылку на коллекцию, создаем фильтр и удаляем игровую сессию в БД.
	collection := s.db.Database(dbName).Collection(colName)
	filter := bson.D{{Key: "_id", Value: idObj}}
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
