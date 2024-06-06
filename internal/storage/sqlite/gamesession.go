package sqlite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"petsittersGameServer/internal/storage"
	"strings"
)

// CreateSession - создает в базе данных нового юзера и игровую сессию для него.
func (s *Storage) CreateSession(ctx context.Context, name, email string) (*storage.GameSession, error) {
	const operation = "storage.sqlite.CreateSession"

	// При POST запросе на создание нового игрока проверяем email. Если такой игрок уже сущетсвует,
	// то возвращаем его игровую сессию
	res, err := s.GetSessionByEmail(ctx, email)
	if err == nil {
		return res, storage.ErrUserExists
	}

	var gs storage.GameSession
	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Создаем строку в таблице users с данными игрока и записываем их в структуру сессии
	err = tx.QueryRowContext(ctx, `
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
	err = tx.QueryRowContext(ctx, `
		INSERT INTO game_sessions (user_id, modules) VALUES (?, ?) RETURNING 
		id, created_at, updated_at, current_module, completed, any_field_one, any_field_two, minigame;
	`, gs.UserID, s.modules).Scan(
		&gs.SessionID,
		&gs.CreatedAt,
		&gs.UpdatedAt,
		&gs.CurrentModule,
		&gs.Completed,
		&gs.AnyFieldOne,
		&gs.AnyFieldTwo,
		&gs.Minigame,
	)
	// Распознаем ошибки из БД и приводим их в более наглядый вид
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	// Декодируем json модулей в структуру
	err = json.Unmarshal([]byte(s.modules), &gs.Modules)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// GetSessionById - возвращает игровую сессию из БД по ее Id.
func (s *Storage) GetSessionById(ctx context.Context, id int) (*storage.GameSession, error) {
	const operation = "storage.sqlite.GetSessionById"
	var gs storage.GameSession
	var mod string

	// Выполняем запрос и записываем результат в структуру GameSession
	err := s.db.QueryRowContext(ctx, `
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id
		WHERE game_sessions.id = (?);
	`, id).Scan(
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
		&mod,
		&gs.Minigame,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	// Декодируем json модулей в структуру
	err = json.Unmarshal([]byte(mod), &gs.Modules)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// GetSessionByEmail - возвращает игровую сессию из БД по емэйлу ее игрока.
func (s *Storage) GetSessionByEmail(ctx context.Context, email string) (*storage.GameSession, error) {
	const operation = "storage.sqlite.GetSessionByEmail"
	var gs storage.GameSession
	var mod string

	// Выполняем запрос и записываем результат в структуру GameSession
	err := s.db.QueryRowContext(ctx, `
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id
		WHERE users.email = (?);
	`, email).Scan(
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
		&mod,
		&gs.Minigame,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	// Декодируем json модулей в структуру
	err = json.Unmarshal([]byte(mod), &gs.Modules)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// GetSessions - возвращает все игровые сессии из БД.
func (s *Storage) GetSessions(ctx context.Context) ([]storage.GameSession, error) {
	const operation = "storage.sqlite.GetSessions"
	var arr []storage.GameSession

	// Выполняем запрос и записываем результат построчно в слайс
	rows, err := s.db.QueryContext(ctx, `
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id 
		ORDER BY game_sessions.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	for rows.Next() {
		var gs storage.GameSession
		var mod string
		err := rows.Scan(
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
			&mod,
			&gs.Minigame,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
		}
		// Декодируем json модулей в структуру
		err = json.Unmarshal([]byte(mod), &gs.Modules)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
		arr = append(arr, gs)
	}
	// Проверяем, была ли получена хотя бы одна строка из БД
	if len(arr) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionsEmpty)
	}
	return arr, nil
}

// CleanSession - устанавливает изменяемые поля игровой сессии в значение по-умолчанию.
func (s *Storage) CleanSession(ctx context.Context, id int) (*storage.GameSession, error) {
	const operation = "storage.sqlite.CleanSession"
	var gs storage.GameSession
	var mod string

	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Очищаем игровую сессию
	res, err := tx.ExecContext(ctx, `
		UPDATE game_sessions SET 
		updated_at = CURRENT_TIMESTAMP,
		current_module = 1,
		completed = false,
		any_field_one = '',
		any_field_two = '',
		modules = (?),
		minigame = '' 
		WHERE id = (?);
	`, s.modules, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if q, _ := res.RowsAffected(); q == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	// Получаем очищенные данные и записываем в структуру GameSession
	err = tx.QueryRowContext(ctx, `
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id
		WHERE game_sessions.id = (?);
	`, id).Scan(
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
		&mod,
		&gs.Minigame,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	// Декодируем json модулей в структуру
	err = json.Unmarshal([]byte(mod), &gs.Modules)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// UpdateSession - устанавливает изменяемые поля игровой сессии в соответствии с переданными значениями.
func (s *Storage) UpdateSession(ctx context.Context, gs storage.GameSession) error {
	const operation = "storage.sqlite.UpdateSession"

	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Проверка наличия модуля.
	// Вернет ошибку, если попытаться вставить в current_module номер несуществующего модуля
	count := strings.Count(s.modules, "module_name")
	if gs.CurrentModule < 1 || gs.CurrentModule > count {
		return fmt.Errorf("%s: %w", operation, storage.ErrModuleNotFound)
	}

	// Подготавливаем строку modules для записи в БД
	mod, err := json.Marshal(gs.Modules)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	buf := new(bytes.Buffer)
	err = json.Compact(buf, mod)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем запрос, вставляя значения из полей полученной структуры
	res, err := tx.ExecContext(ctx, `
		UPDATE game_sessions SET
		updated_at = CURRENT_TIMESTAMP,
		current_module = (?),
		completed = (?),
		any_field_one = (?),
		any_field_two = (?),
		modules = (?),
		minigame = (?)
		WHERE id = (?);
	`,
		gs.CurrentModule,
		gs.Completed,
		gs.AnyFieldOne,
		gs.AnyFieldTwo,
		buf.String(),
		gs.Minigame,
		gs.SessionID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if q, _ := res.RowsAffected(); q == 0 {
		return fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

// DeleteSessionById - удаляет данные об игроке и игровой сессии по ее id.
func (s *Storage) DeleteSessionById(ctx context.Context, id int) error {
	const operation = "storage.sqlite.DeleteSessionById"

	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем удаление игровой сессии
	res, err := tx.ExecContext(ctx, `
		DELETE FROM game_sessions WHERE id = (?);
	`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if q, _ := res.RowsAffected(); q == 0 {
		return fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	// Выполняем удаление игрока
	_, err = tx.ExecContext(ctx, `
		DELETE FROM users WHERE id = (?);
	`, id)
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

// TruncateTables - удаляет все записи из таблиц users и game_sessions.
func (s *Storage) TruncateTables(ctx context.Context) error {
	const operation = "storage.sqlite.TruncateTables"

	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Очищаем таблицу игровых сессий
	_, err = tx.ExecContext(ctx, `
		DELETE FROM game_sessions;
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	// Очищаем таблицу игроков
	_, err = tx.ExecContext(ctx, `
		DELETE FROM users;
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
