package repository

import (
	"context"
	"fmt"

	"github.com/amiosamu/gofemart/internal/domain"
)

func (s *Storage) AddOrder(ctx context.Context, order domain.Order) error {
	result, err := s.DB.ExecContext(ctx, "INSERT INTO orders (order_id, status, uploaded_at, bonuses, user_id) values ($1, $2, $3, $4, $5) on conflict (order_id) do nothing",
		order.OrderID, order.Status, order.UploadedAt, order.Bonuses, order.UserID)
	if err != nil {
		return fmt.Errorf("postgreSQL: addOrder %s", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgreSQL: addOrder %s", err)
	}

	if rowsAffected == 0 {
		userID, err := s.checkOrder(ctx, order)
		if err != nil {
			return fmt.Errorf("postgreSQL: addOrder %s", err)
		}

		if userID == order.UserID {
			return domain.ErrAlreadyUploadedByThisUser
		} else {
			return domain.ErrAlreadyUploadedByAnotherUser
		}
	}

	return nil
}

func (s *Storage) checkOrder(ctx context.Context, order domain.Order) (int64, error) {
	var userID int64
	err := s.DB.QueryRowContext(ctx, "SELECT user_id FROM orders WHERE order_id=$1", order.OrderID).
		Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("postgreSQL: checkOrder %s", err)
	}
	return userID, nil
}

func (s *Storage) GetAllOrders(ctx context.Context, userID int64) ([]domain.Order, error) {
	var orders []domain.Order
	rows, err := s.DB.QueryContext(ctx, "SELECT order_id, status, uploaded_at, bonuses FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("postgreSQL: getAllOrders %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order domain.Order
		err := rows.Scan(&order.OrderID, &order.Status, &order.UploadedAt, &order.Bonuses)
		if err != nil {
			return nil, fmt.Errorf("postgreSQL: getAllOrders %s", err)
		}
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("postgreSQL: getAllOrders %s", err)
	}

	if len(orders) == 0 {
		return nil, domain.ErrNoData
	}

	return orders, nil
}
