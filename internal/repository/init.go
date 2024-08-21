package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var (
	ErrDuplicate = errors.New("login already in use")
)

type Storage struct {
	DB *sql.DB
}


func NewStorage(addr string) (*Storage, error) {
	db, err := goose.OpenDBWithDriver("pgx", addr)
	if err != nil {
		return nil, fmt.Errorf("goose: failed to open DB: %v", err)
	}

	err = goose.Up(db, "./migrations")
	if err != nil {
		return nil, fmt.Errorf("goose: failed to migrate: %v", err)
	}

	return &Storage{
		DB: db,
	}, nil
}


func (s *Storage) Close() error {
	return s.DB.Close()
}
