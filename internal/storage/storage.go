package storage

import "errors"

// Ошибки при работе с БД
var (
	ErrUserExists      = errors.New("user exists")
	ErrInput           = errors.New("incorrect user data input")
	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
)

// GameSession - структура игровой сессии.
type GameSession struct {
	SessionID int `json:"id"`
	User
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	CurrentModule int    `json:"currentScene"`
	Completed     bool   `json:"complete"`
	AnyFieldOne   string `json:"anyFieldOne"`
	AnyFieldTwo   string `json:"anyFieldTwo"`
	Modules       string `json:"modules"`
	Minigame      string `json:"minigame"`
}

// User - структура пользователя.
type User struct {
	UserID   int    `json:"userId"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type Modules struct {
}
