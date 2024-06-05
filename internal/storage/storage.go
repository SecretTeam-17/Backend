package storage

import "errors"

// Ошибки при работе с БД
var (
	ErrUserExists      = errors.New("user exists")
	ErrInput           = errors.New("incorrect user data input")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionsEmpty   = errors.New("table game_sessions is empty")
	ErrUserNotFound    = errors.New("user not found")
	ErrModuleNotFound  = errors.New("module not found")
)

// GameSession - структура игровой сессии.
type GameSession struct {
	SessionID int `json:"id" validate:"required"`
	User
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
	CurrentModule int      `json:"currentModule" validate:"required"`
	Completed     bool     `json:"completed"`
	AnyFieldOne   string   `json:"anyFieldOne"`
	AnyFieldTwo   string   `json:"anyFieldTwo"`
	Modules       []Module `json:"modules"`
	Minigame      string   `json:"minigame"`
}

// User - структура пользователя.
type User struct {
	UserID   int    `json:"userId"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type Module struct {
	Name      string     `json:"name"`
	Available bool       `json:"available"`
	Questions []Question `json:"questions"`
}

type Question struct {
	QuestionId    int      `json:"questionId"`
	QuestionText  string   `json:"questionText"`
	Answers       []Answer `json:"answers"`
	Corrects      []int    `json:"corrects"`
	PlayerAnswers []int    `json:"playerAnswers"`
}

type Answer struct {
	AnswerNum  int    `json:"answerNum"`
	AnswerText string `json:"answerText"`
}
