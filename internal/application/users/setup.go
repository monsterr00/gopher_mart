package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/monsterr00/gopher_mart/internal/application/users/handlers"
)

// func Setup(router *chi.Mux, repo handlers.UserRepository, uCrRepo handlers.UserCreationService) {
func Setup(router *chi.Mux, uCrRepo handlers.UserCreationService) {
	handler := handlers.NewHandler(uCrRepo)
	//handler := handlers.NewHandler(repo, uCrRepo)

	router.Post("/api/user/register", uCrRepo.BasicAuth(handler.Register)) // Регистрация пользователя
	router.Post("/api/user/login", uCrRepo.BasicAuth(handler.Login))       // Аутентификация пользователя
}
