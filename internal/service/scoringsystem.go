package service

import (
	"context"

	"github.com/amiosamu/gofemart/internal/domain"
)

type ScoringSystemRepository interface {
	GetOrderStatus(ctx context.Context) ([]string, error)
	UpdateOrder(ctx context.Context, order domain.ScoringSystem) error
}

type ScoringSystem struct {
	repo ScoringSystemRepository
}

func NewScoringSystem(repo ScoringSystemRepository) *ScoringSystem {
	return &ScoringSystem{
		repo: repo,
	}
}

func (s *ScoringSystem) GetOrderStatus(ctx context.Context) ([]string, error) {
	return s.repo.GetOrderStatus(ctx)
}

func (s *ScoringSystem) UpdateOrder(ctx context.Context, order domain.ScoringSystem) error {
	return s.repo.UpdateOrder(ctx, order)
}
