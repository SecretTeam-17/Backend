package sqlite

import (
	"fmt"
	"petsittersGameServer/internal/storage"

	"github.com/mattn/go-sqlite3"
)

// CreateSession - создает в базе данных нового юзера и игровую сессию для него.
func (s *Storage) CreateSession(name, email string) (*storage.GameSession, error) {
	const operation = "storage.sqlite.CreateSession"
	var gs storage.GameSession

	// TODO: попробовать Prepare с транзакцией
	// Начинаем транзакцию
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Создаем строку в таблице users с данными игрока и записываем их в структуру сессии
	err = tx.QueryRow(`
		INSERT INTO users (name, email) VALUES (?, ?) RETURNING id, name, email;
		`,
		checkNullString(name),
		checkNullString(email),
	).Scan(&gs.UserID, &gs.UserName, &gs.Email)
	// Распознаем ошибки из БД и приводим их в более наглядый вид
	// TODO: вынести в отдельную функцию
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", operation, storage.ErrUserExists)
		}
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintNotNull {
			return nil, fmt.Errorf("%s: %w", operation, storage.ErrInput)
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Создаем строку в таблице game_sessions, вставляем туда id игрока и возвращаем все атрибуты.
	// Записываем их в структуру
	err = tx.QueryRow(`
		INSERT INTO game_sessions (user_id) VALUES (?) RETURNING 
		id, created_at, updated_at, current_module, completed, any_field_one, any_field_two, modules, minigame;
	`, gs.UserID).Scan(
		&gs.SessionID,
		&gs.CreatedAt,
		&gs.UpdatedAt,
		&gs.CurrentModule,
		&gs.Completed,
		&gs.AnyFieldOne,
		&gs.AnyFieldTwo,
		&gs.Modules,
		&gs.Minigame,
	)
	// Распознаем ошибки из БД и приводим их в более наглядый вид
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", operation, storage.ErrUserExists)
		}
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintNotNull {
			return nil, fmt.Errorf("%s: %w", operation, storage.ErrInput)
		}
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}
