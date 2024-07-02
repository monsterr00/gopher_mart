package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/monsterr00/gopher_mart/internal/domain/entities"
	"golang.org/x/crypto/bcrypt"
)

func (ucs *UserCreationService) SaveUser(ctx context.Context, user entities.User) error {
	var err error

	user.Id, _ = uuid.NewRandom()

	user.Password, err = user.HashPassword()
	if err != nil {
		return err
	}

	err = ucs.userRepo.CheckDuplicateLogin(ctx, user.Login)
	if err != nil {
		return err
	}

	err = ucs.userRepo.SaveUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (ucs *UserCreationService) CheckAuth(ctx context.Context, user entities.User) error {
	passwordHash, err := ucs.userRepo.GetPassword(ctx, user.Login)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	if err != nil {
		return err
	}

	return nil
}
