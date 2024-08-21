package domain

import (
	"errors"
)

var (
	ErrNoWithdraws = errors.New("the user has no withdraws")
	ErrNoBonuses   = errors.New("not enough bonuses")
)

type Withdraw struct {
	OrderID    string  `json:"order"`
	Bonuses    float32 `json:"sum"`
	UploadedAt string  `json:"processed_at"`
	UserID     int64   `json:"-"`
}

type BalanceOutput struct {
	Bonuses  float32 `json:"current"`
	Withdraw float32 `json:"withdrawn"`
}
