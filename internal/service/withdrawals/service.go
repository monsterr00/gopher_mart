package withdrawals

import (
	"context"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

type WithdrawalRepository interface {
	GetBalanceAdded(ctx context.Context, login string) (float64, error)
	GetBalanceSpended(ctx context.Context, login string) (float64, error)
	GetWithdrawals(ctx context.Context, login string) ([]entities.Withdrawal, error)
	SaveWithdrawal(ctx context.Context, order entities.Order) error
}

type OrderRepository interface {
	SaveOrder(ctx context.Context, order entities.Order) error
	CheckOrderDuplicate(ctx context.Context, order entities.Order) error
	UpdateStatus(ctx context.Context, order entities.RegistredOrder) error
	UpdateBonuses(ctx context.Context, order entities.RegistredOrder) error
	GetOrdersLogin(ctx context.Context, login string) ([]entities.Order, error)
}

type WithdrawalCreationService struct {
	orderRepo      OrderRepository
	withdrawalRepo WithdrawalRepository
}

func NewWithdrawalCreationService(or OrderRepository, wr WithdrawalRepository) *WithdrawalCreationService {
	return &WithdrawalCreationService{
		orderRepo:      or,
		withdrawalRepo: wr,
	}
}
