package sqlite

import "database/sql"

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
