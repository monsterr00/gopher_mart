package handlers

import (
	"context"
	"net/http"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

type UserCreationService interface {
	BasicAuth(next http.HandlerFunc) http.HandlerFunc
	SaveUser(ctx context.Context, user entities.User) error
	CheckAuth(ctx context.Context, user entities.User) error
}

type WithdrawalCreationService interface {
	GetBalance(ctx context.Context, login string) (entities.Balance, error)
	GetWithdrawals(ctx context.Context, login string) ([]entities.FormatedWithdrawal, error)
	SaveOrder(ctx context.Context, order entities.Order) error
}

type WithdrawalHandler struct {
	uCrRepo UserCreationService
	wCrRepo WithdrawalCreationService
}

func NewHandler(uCrR UserCreationService, wCrR WithdrawalCreationService) *WithdrawalHandler {
	return &WithdrawalHandler{
		uCrRepo: uCrR,
		wCrRepo: wCrR,
	}
}
