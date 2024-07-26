package withdrawals

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/monsterr00/gopher_mart/internal/domain/entities"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

func (wcs *WithdrawalCreationService) SaveOrder(ctx context.Context, order entities.Order) error {
	err := wcs.orderRepo.CheckOrderDuplicate(ctx, order)
	if err != nil && !errors.Is(err, i.ErrDuplicateUserOrder) {
		return err
	}

	if err == nil {
		order.ID, _ = uuid.NewRandom()
		err = wcs.orderRepo.SaveOrder(ctx, order)
		if err != nil {
			return err
		}
	}

	balance, err := wcs.GetBalance(ctx, order.Login)
	if err != nil {
		return err
	}

	if balance.Balance < order.SpendedBonuses {
		return i.ErrInsufFunds
	}

	err = wcs.withdrawalRepo.SaveWithdrawal(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
