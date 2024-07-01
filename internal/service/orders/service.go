package orders

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

type UserRepository interface {
	GetUserByLogin(ctx context.Context, login string) (entities.User, error)
	SaveUser(ctx context.Context, user entities.User) error
	GetPassword(ctx context.Context, login string) (string, error)
	CheckDuplicateLogin(ctx context.Context, login string) error
}

type OrderRepository interface {
	SaveOrder(ctx context.Context, order entities.Order) error
	CheckOrderDuplicate(ctx context.Context, order entities.Order) error
	UpdateStatus(ctx context.Context, order entities.RegistredOrder) error
	UpdateBonuses(ctx context.Context, order entities.RegistredOrder) error
	GetOrdersLogin(ctx context.Context, login string) ([]entities.Order, error)
}

type OrderCreationService struct {
	userRepo  UserRepository
	orderRepo OrderRepository
	host      string
	client    *resty.Client
}

func NewOrderCreationService(or OrderRepository, ur UserRepository, host string, client *resty.Client) *OrderCreationService {
	return &OrderCreationService{
		orderRepo: or,
		userRepo:  ur,
		host:      host,
		client:    client,
	}
}
