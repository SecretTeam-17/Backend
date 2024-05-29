package sqlite

import (
	"database/sql"
	"errors"
	"petsittersGameServer/internal/storage"

	"github.com/mattn/go-sqlite3"
)

// checkNullString - проверяет строку и если она пустая, то преобразует ее в NULL значение для передачи в БД.
func checkNullString(str string) sql.NullString {
	if len(str) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: str,
		Valid:  true,
	}
}

// checkDBError - выявляет конкретный тип ошибки из БД.
func checkDBError(err error) error {
	sqliteErr, ok := err.(sqlite3.Error)
	if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return storage.ErrUserExists
	}
	if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintNotNull {
		return storage.ErrInput
	}

	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrSessionNotFound
	}
	return err
}
