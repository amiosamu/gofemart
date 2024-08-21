package repository

import (
	"context"
	"fmt"

	"github.com/amiosamu/gofemart/internal/domain"
)


func (s *Storage) Create(ctx context.Context, user domain.User) error {
	result, err := s.DB.ExecContext(ctx, "INSERT INTO users (login, password, registered_at) values ($1, $2, $3) on conflict (login) do nothing",
		user.Login, user.Password, user.RegisteredAt)
	if err != nil {
		return fmt.Errorf("postgreSQL: create %s", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgreSQL: create %s", err)
	}

	if rowsAffected == 0 {
		return ErrDuplicate
	}

	return nil
}


func (s *Storage) GetUser(ctx context.Context, login, password string) (domain.User, error) {
	var user domain.User
	err := s.DB.QueryRowContext(ctx, "SELECT id, login, password, registered_at FROM users WHERE login=$1 AND password=$2", login, password).
		Scan(&user.ID, &user.Login, &user.Password, &user.RegisteredAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("postgreSQL: getUser %s", err)
	}
	return user, nil
}
