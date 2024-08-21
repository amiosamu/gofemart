package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/amiosamu/gofemart/internal/domain"
)

type OrderRepository interface {
	AddOrder(ctx context.Context, order domain.Order) error
	GetAllOrders(ctx context.Context, userID int64) ([]domain.Order, error)
}

type Orders struct {
	repo OrderRepository
}

func NewOrders(repo OrderRepository) *Orders {
	return &Orders{
		repo: repo,
	}
}

// AddOrderID загружает номер заказа в систему.
func (o *Orders) AddOrderID(ctx context.Context, orderID string) error {
	trimmedStr := strings.TrimSpace(orderID)
	if len(trimmedStr) == 0 {
		return domain.ErrIncorrectOrder
	}

	if !checkOrderNumber(orderID) {
		return domain.ErrIncorrectOrder
	}

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return errors.New("incorrect user id")
	}

	order := domain.Order{
		OrderID:    orderID,
		Status:     domain.NewOrder,
		UploadedAt: time.Now().Format(time.RFC3339),
		Bonuses:    0,
		UserID:     userID,
	}

	return o.repo.AddOrder(ctx, order)
}

// GetAllOrders выводит отсортированный по дате список заказов пользователя.
func (o *Orders) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	orders, err := o.repo.GetAllOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
