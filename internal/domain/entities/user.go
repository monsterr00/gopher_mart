package entities

import (
	"time"

	"github.com/google/uuid"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uuid.UUID
	Login          string    `json:"login"`    // Логин
	Password       string    `json:"password"` // Пароль
	CreatedAt      time.Time // Дата регистрации
	Balance        float64   // Баланс
	SpendedBonuses float64   // Сумма списанных бнусов
}

func (u User) ValidateUser() error {
	if !u.validateLogin() || !u.validatePassword() {
		return i.ErrValidate
	}
	return nil
}

func (u User) validateLogin() bool {
	return u.Login != ""
}

func (u User) validatePassword() bool {
	return u.Password != ""
}

func (u User) HashPassword() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
