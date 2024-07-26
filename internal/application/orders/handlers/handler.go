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

type OrderCreationService interface {
	SaveOrder(ctx context.Context, order entities.Order) error
	GetOrdersLogin(ctx context.Context, login string) ([]entities.FormatedOrder, error)
}

type OrderHandlers struct {
	uCrRepo UserCreationService
	oCrRepo OrderCreationService
}

func NewHandler(uCrR UserCreationService, oCrR OrderCreationService) *OrderHandlers {
	return &OrderHandlers{
		uCrRepo: uCrR,
		oCrRepo: oCrR,
	}
}
