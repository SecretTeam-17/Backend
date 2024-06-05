package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Таймаут на создание талблиц в БД
const timeout time.Duration = time.Second * 10

// Storage - пул подключений к БД
type Storage struct {
	db      *sql.DB
	modules string
}

// New - конструктор пула подключений к БД. Создает также пустые таблицы, если их нет.
func New(storagePath, questionPath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	// Создаем пул и проверяем его
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = initTables(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	mod, err := initModules(questionPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db, modules: mod}, nil
}

func initTables(ctx context.Context, db *sql.DB) error {
	const operation = "storage.sqlite.initTables"

	// Начинаем транзакцию
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Создаем таблицы:
	// Таблица users
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE CHECK(email LIKE '%@%')
		);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	// Индекс по полю email в таблице users
	_, err = tx.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_email ON users(email);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	// Таблица game_sessions
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS game_sessions (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			current_module INTEGER REFERENCES modules(id) DEFAULT 1,
			completed BOOLEAN DEFAULT FALSE,
			any_field_one TEXT DEFAULT '',
			any_field_two TEXT DEFAULT '',
			modules TEXT NOT NULL DEFAULT '',
			minigame TEXT NOT NULL DEFAULT ''
		);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

// Close - обертка для закрытия подключения к БД.
func (s *Storage) Close() error {
	return s.db.Close()
}

func initModules(path string) (string, error) {
	const operation = "storage.sqlite.initModules"

	// Читаем из json файла в строку
	file, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	// Вырезаем лишние пробелы и отступы
	buf := new(bytes.Buffer)
	err = json.Compact(buf, file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	return buf.String(), nil
}
