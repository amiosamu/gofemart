package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/amiosamu/gofemart/internal/domain"
	"github.com/amiosamu/gofemart/internal/repository"
	"github.com/shopspring/decimal"
)

type BonusesRepository interface {
	Balance(ctx context.Context, userID int64) (float32, error)
	WithdrawBalance(ctx context.Context, userID int64) (float32, error)
	Withdraw(ctx context.Context, withdraw domain.Withdraw) error
	Withdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error)
}

type Bonuses struct {
	repo    BonusesRepository
	storage *repository.Storage
}

func NewBonuses(repo BonusesRepository, storage *repository.Storage) *Bonuses {
	return &Bonuses{
		repo:    repo,
		storage: storage,
	}
}

func (b *Bonuses) Balance(ctx context.Context) (*domain.BalanceOutput, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	balanceUser, err := b.repo.Balance(ctx, userID)
	if err != nil {
		return nil, err
	}

	balanceWithdraws, err := b.repo.WithdrawBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	newBalanceUser := decimal.NewFromFloat32(balanceUser).Sub(decimal.NewFromFloat32(balanceWithdraws))

	var balance domain.BalanceOutput
	balance.Bonuses = float32(newBalanceUser.InexactFloat64())
	balance.Withdraw = balanceWithdraws
	return &balance, nil
}

func (b *Bonuses) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {
	trimmedStr := strings.TrimSpace(withdraw.OrderID)
	if len(trimmedStr) == 0 {
		return domain.ErrIncorrectOrder
	}

	if !checkOrderNumber(withdraw.OrderID) {
		return domain.ErrIncorrectOrder
	}

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return errors.New("incorrect user id")
	}

	with := domain.Withdraw{
		OrderID:    withdraw.OrderID,
		Bonuses:    withdraw.Bonuses,
		UploadedAt: time.Now().Format(time.RFC3339),
		UserID:     userID,
	}

	// Начало транзакции
	tx, err := b.storage.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = b.repo.Withdraw(ctx, with)
	if err != nil {
		return err
	}

	balanceUser, err := b.repo.Balance(ctx, userID)
	if err != nil {
		return err
	}

	balanceWithdraws, err := b.repo.WithdrawBalance(ctx, userID)
	if err != nil {
		return err
	}

	sum := decimal.NewFromFloat32(balanceUser).Sub(decimal.NewFromFloat32(balanceWithdraws))
	if sum.LessThan(decimal.Zero) {
		tx.Rollback()
		return domain.ErrNoBonuses
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *Bonuses) Withdrawals(ctx context.Context) ([]domain.Withdraw, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	withdrawals, err := b.repo.Withdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) <= 10 {
		return withdrawals, nil
	}

	startIndex := len(withdrawals) - 10
	endIndex := len(withdrawals)
	return withdrawals[startIndex:endIndex], nil
}
