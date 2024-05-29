package sqlite

import (
	"fmt"
	"petsittersGameServer/internal/storage"
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
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
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
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// GetSessionById - возвращает игровую сессию из БД по ее Id.
func (s *Storage) GetSessionById(id int) (*storage.GameSession, error) {
	const operation = "storage.sqlite.GetSessionById"
	var gs storage.GameSession

	// Подготавливаем запрос
	stmt, err := s.db.Prepare(`
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id
		WHERE users.id = (?);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем запрос и записываем результат в структуру GameSession
	err = stmt.QueryRow(id).Scan(
		&gs.SessionID,
		&gs.UserID,
		&gs.UserName,
		&gs.Email,
		&gs.CreatedAt,
		&gs.UpdatedAt,
		&gs.CurrentModule,
		&gs.Completed,
		&gs.AnyFieldOne,
		&gs.AnyFieldTwo,
		&gs.Modules,
		&gs.Minigame,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return &gs, nil
}

// GetSessionByEmail - возвращает игровую сессию из БД по емэйлу ее игрока.
func (s *Storage) GetSessionByEmail(email string) (*storage.GameSession, error) {
	const operation = "storage.sqlite.GetSessionByEmail"
	var gs storage.GameSession

	// Подготавливаем запрос
	stmt, err := s.db.Prepare(`
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id
		WHERE users.email = (?);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем запрос и записываем результат в структуру GameSession
	err = stmt.QueryRow(email).Scan(
		&gs.SessionID,
		&gs.UserID,
		&gs.UserName,
		&gs.Email,
		&gs.CreatedAt,
		&gs.UpdatedAt,
		&gs.CurrentModule,
		&gs.Completed,
		&gs.AnyFieldOne,
		&gs.AnyFieldTwo,
		&gs.Modules,
		&gs.Minigame,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}

	return &gs, nil
}
