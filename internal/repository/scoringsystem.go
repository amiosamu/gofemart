package repository

import (
	"context"
	"fmt"

	"github.com/amiosamu/gofemart/internal/domain"
)

func (s *Storage) GetOrderStatus(ctx context.Context) ([]string, error) {
	var orderID []string
	rows, err := s.DB.QueryContext(ctx, "SELECT order_id FROM orders WHERE status NOT IN ('PROCESSED', 'INVALID') LIMIT 15")
	if err != nil {
		return nil, fmt.Errorf("postgreSQL: getOrderStatus %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("postgreSQL: getOrderStatus %s", err)
		}
		orderID = append(orderID, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgreSQL: getOrderStatus %s", err)
	}

	return orderID, nil
}

func (s *Storage) UpdateOrder(ctx context.Context, order domain.ScoringSystem) error {
	_, err := s.DB.ExecContext(ctx, "UPDATE orders SET status=$1, bonuses=$2 WHERE order_id=$3", order.Status, order.Bonuses, order.OrderID)
	if err != nil {
		return fmt.Errorf("postgreSQL: updateOrder %s", err)
	}
	return nil
}
