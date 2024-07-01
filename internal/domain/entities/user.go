package entities

import (
	"time"

	"github.com/google/uuid"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

// TO-DO
/*
 Добавить валидацию полей логин и пароль
*/

type User struct {
	Id             uuid.UUID
	Login          string    `json:"login"`    // Логин
	Password       string    `json:"password"` // Пароль
	CreatedAt      time.Time // Дата регистрации
	Balance        float64   // Баланс
	SpendedBonuses float64   // Сумма списанных бнусов
}

func NewUser(id uuid.UUID, login string, password string) *User {
	return &User{
		Id:             id,
		Login:          login,
		Password:       password,
		CreatedAt:      time.Now(),
		Balance:        0,
		SpendedBonuses: 0,
	}
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
