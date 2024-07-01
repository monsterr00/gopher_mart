package orders

import (
	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/gopher_mart/internal/application/orders/handlers"
)

func Setup(router *chi.Mux, uCrRepo handlers.UserCreationService, oCrRepo handlers.OrderCreationService) {
	handler := handlers.NewHandler(uCrRepo, oCrRepo)

	router.Post("/api/user/orders", uCrRepo.BasicAuth(handler.Create)) // Загрузка пользователем номера заказа для расчёта
	router.Get("/api/user/orders", uCrRepo.BasicAuth(handler.Fetch))   // Получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
}
