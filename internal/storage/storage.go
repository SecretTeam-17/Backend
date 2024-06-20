package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ошибки при работе с БД
var (
	ErrUserExists      = errors.New("user exists")
	ErrInput           = errors.New("incorrect user data input")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionsEmpty   = errors.New("collection gamesessions is empty")
	ErrUserNotFound    = errors.New("user not found")
	ErrModuleNotFound  = errors.New("module not found")
)

type GameSession struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
	Stats     json.RawMessage    `json:"stats" bson:"stats"`
	Modules   json.RawMessage    `json:"modules" bson:"modules"`
	Minigames json.RawMessage    `json:"minigames" bson:"minigames"`
}

type Interface interface {
	CleanSession(ctx context.Context, id string) (*GameSession, error)
	CreateSession(ctx context.Context, name, email string, stats, modules, minigames []byte) (*GameSession, error)
	DeleteSessionById(ctx context.Context, id string) error
	GetSessions(ctx context.Context) ([]GameSession, error)
	GetSessionByEmail(ctx context.Context, email string) (*GameSession, error)
	GetSessionById(ctx context.Context, id string) (*GameSession, error)
	TruncateData(ctx context.Context) error
	UpdateSession(ctx context.Context, gs GameSession) error
}
