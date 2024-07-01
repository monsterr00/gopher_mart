package users

import (
	"context"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

type UserRepository interface {
	GetUserByLogin(ctx context.Context, login string) (entities.User, error)
	SaveUser(ctx context.Context, user entities.User) error
	GetPassword(ctx context.Context, login string) (string, error)
	CheckDuplicateLogin(ctx context.Context, login string) error
}

type UserCreationService struct {
	userRepo UserRepository
}

func NewUserCreationService(ur UserRepository) *UserCreationService {
	return &UserCreationService{
		userRepo: ur,
	}
}
