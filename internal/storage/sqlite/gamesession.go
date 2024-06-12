package sqlite

import (
	"context"
	"fmt"
	"petsittersGameServer/internal/storage"
)

// Количество модулей игры зашито в константу пока не ясно, как получать это число динамически.
const modulesCount int = 4

// CreateSession - создает в базе данных нового юзера и игровую сессию для него.
func (s *Storage) CreateSession(ctx context.Context, name, email string) (*GameSession, error) {
	const operation = "storage.sqlite.CreateSession"

	// При POST запросе на создание нового игрока проверяем email. Если такой игрок уже сущетсвует,
	// то возвращаем его игровую сессию
	res, err := s.GetSessionByEmail(ctx, email)
	if err == nil {
		return res, storage.ErrUserExists
	}

	var gs GameSession
	var mod string
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
		&mod,
		&gs.Minigame,
	)
	// Распознаем ошибки из БД и приводим их в более наглядый вид
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	gs.Modules.UnmarshalJSON([]byte(mod))

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// GetSessionById - возвращает игровую сессию из БД по ее Id.
func (s *Storage) GetSessionById(ctx context.Context, id int) (*GameSession, error) {
	const operation = "storage.sqlite.GetSessionById"
	var gs GameSession
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
	gs.Modules.UnmarshalJSON([]byte(mod))

	return &gs, nil
}

// GetSessionByEmail - возвращает игровую сессию из БД по емэйлу ее игрока.
func (s *Storage) GetSessionByEmail(ctx context.Context, email string) (*GameSession, error) {
	const operation = "storage.sqlite.GetSessionByEmail"
	var gs GameSession
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
	gs.Modules.UnmarshalJSON([]byte(mod))

	return &gs, nil
}

// GetSessions - возвращает все игровые сессии из БД.
func (s *Storage) GetSessions(ctx context.Context) ([]GameSession, error) {
	const operation = "storage.sqlite.GetSessions"
	var arr []GameSession

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
		var gs GameSession
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
		gs.Modules.UnmarshalJSON([]byte(mod))
		arr = append(arr, gs)
	}
	// Проверяем, была ли получена хотя бы одна строка из БД
	if len(arr) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionsEmpty)
	}
	return arr, nil
}

// CleanSession - устанавливает изменяемые поля игровой сессии в значение по-умолчанию.
func (s *Storage) CleanSession(ctx context.Context, id int) (*GameSession, error) {
	const operation = "storage.sqlite.CleanSession"
	var gs GameSession
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
		modules = '',
		minigame = '' 
		WHERE id = (?);
	`, id)
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
	gs.Modules.UnmarshalJSON([]byte(mod))

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// UpdateSession - устанавливает изменяемые поля игровой сессии в соответствии с переданными значениями.
func (s *Storage) UpdateSession(ctx context.Context, gs GameSession) error {
	const operation = "storage.sqlite.UpdateSession"

	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Проверка наличия модуля.
	// Вернет ошибку, если попытаться вставить в current_module номер несуществующего модуля
	if gs.CurrentModule < 1 || gs.CurrentModule > modulesCount {
		return fmt.Errorf("%s: %w", operation, storage.ErrModuleNotFound)
	}

	var mod []byte
	mod, err = gs.Modules.MarshalJSON()
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
		string(mod),
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
func (s *Storage) TruncateData(ctx context.Context) error {
	const operation = "storage.sqlite.TruncateData"

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
