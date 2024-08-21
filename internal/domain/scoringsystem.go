package domain

type ScoringSystem struct {
	OrderID string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Bonuses float32     `json:"accrual"`
}
