package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Storage - пул подключений к БД
type Storage struct {
	db *sql.DB
}

// New - конструктор пула подключений к БД. Создает также пустые таблицы, если их нет.
func New(storagePath string) (*Storage, error) {
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

	// Начинаем транзакцию
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Создаем таблицы:
	// Таблица modules
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS modules (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT DEFAULT ''
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Таблица questions
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS questions (
			id INTEGER PRIMARY KEY,
			module INTEGER NOT NULL REFERENCES modules(id),
			text TEXT NOT NULL,
			has_input BOOLEAN DEFAULT FALSE
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Таблица answers
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS answers (
			id INTEGER PRIMARY KEY,
			question_id INTEGER NOT NULL REFERENCES questions(id),
			text TEXT NOT NULL,
			is_correct BOOLEAN DEFAULT FALSE
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Таблица users
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE CHECK(email LIKE '%@%')
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Индекс по полю email в таблице users
	_, err = tx.Exec(`
		CREATE INDEX IF NOT EXISTS idx_email ON users(email);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Таблица game_sessions
	_, err = tx.Exec(`
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
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Вносим обязательные данные в таблицу module
	// Пока 4 модуля, потом заменим на реальные данные
	_, err = tx.Exec(`
		INSERT OR IGNORE INTO modules (name, description) VALUES
		('module_name_1', 'description 1'),
		('module_name_2', 'description 2'),
		('module_name_3', 'description 3'),
		('module_name_4', 'description 4');
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

// Close - обертка для закрытия подключения к БД.
func (s *Storage) Close() error {
	return s.db.Close()
}
