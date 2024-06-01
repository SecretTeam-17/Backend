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
		WHERE game_sessions.id = (?);
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

// GetSessions - возвращает все игровые сессии из БД.
func (s *Storage) GetSessions() ([]storage.GameSession, error) {
	const operation = "storage.sqlite.GetSessions"
	var arr []storage.GameSession

	// Подготавливаем запрос
	stmt, err := s.db.Prepare(`
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id 
		ORDER BY game_sessions.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем запрос и записываем результат построчно в слайс
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
	}
	for rows.Next() {
		var gs storage.GameSession
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
			&gs.Modules,
			&gs.Minigame,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, checkDBError(err))
		}
		arr = append(arr, gs)
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionsEmpty)
	}
	return arr, nil
}

// CleanSession - устанавливает изменяемые поля игровой сессии в значение по-умолчанию.
func (s *Storage) CleanSession(id int) (*storage.GameSession, error) {
	const operation = "storage.sqlite.CleanSession"
	var gs storage.GameSession

	// Подготавливаем запрос на очистку данных
	cleanStmt, err := s.db.Prepare(`
		UPDATE game_sessions SET 
		updated_at = CURRENT_TIMESTAMP,
		current_module = 1,
		completed = false,
		any_field_one = '',
		any_field_two = '',
		modules = '',
		minigame = '' 
		WHERE id = (?);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	// Подготавливаем запрос на получение очищенных данных
	getStmt, err := s.db.Prepare(`
		SELECT game_sessions.id, users.id AS userId, users.name, users.email, 
		game_sessions.created_at, game_sessions.updated_at, game_sessions.current_module, 
		game_sessions.completed, game_sessions.any_field_one, game_sessions.any_field_two, 
		game_sessions.modules, game_sessions.minigame
		FROM game_sessions
		JOIN users ON users.id = game_sessions.user_id
		WHERE game_sessions.id = (?);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Очищаем игровую сессию
	res, err := tx.Stmt(cleanStmt).Exec(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if q, _ := res.RowsAffected(); q == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}

	// Получаем очищенные данные и записываем в структуру GameSession
	err = tx.Stmt(getStmt).QueryRow(id).Scan(
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

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &gs, nil
}

// UpdateSession - устанавливает изменяемые поля игровой сессии в соответствии с переданными значениями.
func (s *Storage) UpdateSession(gs storage.GameSession) error {
	const operation = "storage.sqlite.UpdateSession"

	// Костыль, так как в SQLite по дефолту не включена поддержка foreign keys.
	// Подготавливаем запрос на проверку наличия модуля в БД
	checkStmt, err := s.db.Prepare(`
		SELECT * FROM modules WHERE id = (?);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	// Подготавливаем запрос на изменение данных игровой сессии
	updateStmt, err := s.db.Prepare(`
		UPDATE game_sessions SET
		updated_at = CURRENT_TIMESTAMP,
		current_module = (?),
		completed = (?),
		any_field_one = (?),
		any_field_two = (?),
		modules = (?),
		minigame = (?)
		WHERE id = (?);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Проверка наличия модуля
	// Вернет ошибку, если попытаться вставить в current_module несуществующий в БД модуль
	rows, err := tx.Stmt(checkStmt).Query(gs.CurrentModule)
	if !rows.Next() {
		return fmt.Errorf("%s: %w", operation, storage.ErrModuleNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем запрос, вставляя значения из полей полученной структуры
	res, err := tx.Stmt(updateStmt).Exec(
		gs.CurrentModule,
		gs.Completed,
		gs.AnyFieldOne,
		gs.AnyFieldTwo,
		gs.Modules,
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
func (s *Storage) DeleteSessionById(id int) error {
	const operation = "storage.sqlite.DeleteSessionById"

	// Подготавливаем запрос на удаление игрока
	userStmt, err := s.db.Prepare(`
		DELETE FROM users WHERE id = (?);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	// Подготавливаем запрос на удаление игрока
	gsStmt, err := s.db.Prepare(`
		DELETE FROM game_sessions WHERE id = (?);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем удаление игровой сессии
	res, err := tx.Stmt(gsStmt).Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	if q, _ := res.RowsAffected(); q == 0 {
		return fmt.Errorf("%s: %w", operation, storage.ErrSessionNotFound)
	}
	// Выполняем удаление игрока
	_, err = tx.Stmt(userStmt).Exec(id)
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
func (s *Storage) TruncateTables() error {
	const operation = "storage.sqlite.TruncateTables"

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Очищаем таблицу игровых сессий
	_, err = tx.Exec(`
		DELETE FROM game_sessions;
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	// Очищаем таблицу игроков
	_, err = tx.Exec(`
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
