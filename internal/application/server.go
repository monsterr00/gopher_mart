package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/gopher_mart/internal/application/orders"
	o "github.com/monsterr00/gopher_mart/internal/application/orders/handlers"
	"github.com/monsterr00/gopher_mart/internal/application/users"
	"github.com/monsterr00/gopher_mart/internal/application/withdrawals"
	w "github.com/monsterr00/gopher_mart/internal/application/withdrawals/handlers"
)

func NewServer(uCrRepo o.UserCreationService, oCrRepo o.OrderCreationService, wCrRepo w.WithdrawalCreationService) *chi.Mux {
	router := chi.NewRouter()

	users.Setup(router, uCrRepo)
	orders.Setup(router, uCrRepo, oCrRepo)
	withdrawals.Setup(router, uCrRepo, wCrRepo)

	return router
}
