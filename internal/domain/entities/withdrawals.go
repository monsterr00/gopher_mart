package entities

import "time"

type Balance struct {
	Balance        float64 `json:"current"`   // Баланс
	SpendedBonuses float64 `json:"withdrawn"` // Списанные бонусы
}

type Withdrawal struct {
	OrderNum  string    `json:"order"`        // Номер заказа
	Bonuses   float64   `json:"sum"`          // Списанные бонусы
	SpendDate time.Time `json:"processed_at"` // Дата списания
}

type FormatedWithdrawal struct {
	OrderNum     string  `json:"order"`        // Номер заказа
	Bonuses      float64 `json:"sum"`          // Списанные бонусы
	FormatedDate string  `json:"processed_at"` // Дата списания
}
