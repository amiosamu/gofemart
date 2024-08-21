package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/amiosamu/gofemart/internal/domain"
)

func (s *Storage) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	result, err := s.DB.ExecContext(ctx, "INSERT INTO withdrawals (order_id, bonuses, uploaded_at, user_id) values ($1, $2, $3, $4) on conflict (order_id) do nothing",
		withdraw.OrderID, withdraw.Bonuses, withdraw.UploadedAt, withdraw.UserID)
	if err != nil {
		return fmt.Errorf("postgreSQL: withdraw %s", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgreSQL: withdraw %s", err)
	}

	if rowsAffected == 0 {
		userID, err := s.checkWithdraw(ctx, withdraw)
		if err != nil {
			return fmt.Errorf("postgreSQL: withdraw %s", err)
		}

		if userID == withdraw.UserID {
			return domain.ErrAlreadyUploadedByThisUser
		} else {
			return domain.ErrAlreadyUploadedByAnotherUser
		}
	}

	return nil
}

func (s *Storage) Withdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error) {
	var withdrawals []domain.Withdraw
	rows, err := s.DB.QueryContext(ctx, "SELECT order_id, bonuses, uploaded_at FROM withdrawals WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("postgreSQL: withdrawals %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var withdraw domain.Withdraw
		err := rows.Scan(&withdraw.OrderID, &withdraw.Bonuses, &withdraw.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("postgreSQL: withdrawals %s", err)
		}
		withdrawals = append(withdrawals, withdraw)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("postgreSQL: withdrawals %s", err)
	}

	if len(withdrawals) == 0 {
		return nil, domain.ErrNoWithdraws
	}

	return withdrawals, nil
}

func (s *Storage) Balance(ctx context.Context, userID int64) (float32, error) {
	var nullableBalance sql.NullFloat64
	err := s.DB.QueryRowContext(ctx, "SELECT SUM(bonuses) FROM orders WHERE user_id=$1", userID).
		Scan(&nullableBalance)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: balance %s", err)
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, nil
}

func (s *Storage) WithdrawBalance(ctx context.Context, userID int64) (float32, error) {
	var nullableBalance sql.NullFloat64
	err := s.DB.QueryRowContext(ctx, "SELECT SUM(bonuses) FROM withdrawals WHERE user_id=$1", userID).
		Scan(&nullableBalance)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: withdrawBalance %s", err)
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, nil
}

func (s *Storage) checkWithdraw(ctx context.Context, withdraw domain.Withdraw) (int64, error) {
	var userID int64
	err := s.DB.QueryRowContext(ctx, "SELECT user_id FROM withdrawals WHERE order_id=$1", withdraw.OrderID).
		Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: checkWithdraw %s", err)
	}
	return userID, nil
}
