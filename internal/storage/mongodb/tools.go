package mongodb

import (
	"errors"
	"petsittersGameServer/internal/storage"

	"go.mongodb.org/mongo-driver/mongo"
)

// checkDBError - выявляет конкретный тип ошибки из БД MongoDB.
func checkDBError(err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return storage.ErrSessionNotFound
	}
	return err
}
