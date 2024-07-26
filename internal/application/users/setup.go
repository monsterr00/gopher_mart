package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/gopher_mart/internal/application/users/handlers"
)

func Setup(router *chi.Mux, uCrRepo handlers.UserCreationService) {
	handler := handlers.NewHandler(uCrRepo)

	router.Post("/api/user/register", handler.Register) // Регистрация пользователя
	router.Post("/api/user/login", handler.Login)       // Аутентификация пользователя
}
