package entities

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id             uuid.UUID //Id
	OrderNum       string    `json:"number"`            // Номер заказа
	Login          string    `json:"login"`             // Логин
	CreatedAt      time.Time `json:"uploaded_at"`       // Дата регистрации
	Status         string    `json:"status"`            // Статус
	AddedBonuses   *float64  `json:"accrual,omitempty"` // Начисленные бонусы
	SpendedBonuses float64   // Списанные бнусы
	SpendDate      time.Time // Дата списания
}

type RegistredOrder struct {
	OrderNum     string   `json:"order"`             // Номер заказа
	Status       string   `json:"status"`            // Статус
	AddedBonuses *float64 `json:"accrual,omitempty"` // Начисленные бонусы
}

type FormatedOrder struct {
	OrderNum     string   `json:"number"`            // Номер заказа
	Status       string   `json:"status"`            // Статус
	AddedBonuses *float64 `json:"accrual,omitempty"` // Начисленные бонусы
	CreatedAt    string   `json:"uploaded_at"`       // Дата регистрации
}

func (o *Order) ValidateOrderNum() bool {
	orderNum, err := strconv.Atoi(o.OrderNum)
	if o.OrderNum == "" || err != nil {
		return false
	}

	return (orderNum%10+checksum(orderNum/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
