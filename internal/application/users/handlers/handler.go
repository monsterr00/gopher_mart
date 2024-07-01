package handlers

import (
	"context"
	"net/http"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

type UserCreationService interface {
	BasicAuth(next http.HandlerFunc) http.HandlerFunc
	SaveUser(ctx context.Context, user entities.User) error
}

type UserHandler struct {
	uCrRepo UserCreationService
}

func NewHandler(uCrR UserCreationService) *UserHandler {
	return &UserHandler{
		uCrRepo: uCrR,
	}
}
