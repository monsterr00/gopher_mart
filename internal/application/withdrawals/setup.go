package withdrawals

import (
	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/gopher_mart/internal/application/withdrawals/handlers"
)

func Setup(router *chi.Mux, uCrRepo handlers.UserCreationService, wCrRepo handlers.WithdrawalCreationService) {
	handler := handlers.NewHandler(uCrRepo, wCrRepo)

	router.Get("/api/user/balance", uCrRepo.BasicAuth(handler.GetBalance))         // Получение текущего баланса счёта баллов лояльности пользователя
	router.Get("/api/user/withdrawals", uCrRepo.BasicAuth(handler.Fetch))          // Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
	router.Post("/api/user/balance/withdraw", uCrRepo.BasicAuth(handler.Withdraw)) // Запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
}
