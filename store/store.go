package store

import (
	"errors"
)

var (
	ErrTxNotFound = errors.New("TRANSACTION NOT FOUND")
)

// Store provides operations for reading and writing data to the database
type Store interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Delete(key string)
	Count(value string) int
	Begin() Store
	Commit() Store
	Rollback() (Store, error)
}
