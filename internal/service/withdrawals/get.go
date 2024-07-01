package withdrawals

import (
	"context"
	"time"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

func (wcs *WithdrawalCreationService) GetBalance(ctx context.Context, login string) (entities.Balance, error) {
	var balance entities.Balance

	added, err := wcs.withdrawalRepo.GetBalanceAdded(ctx, login)
	if err != nil {
		return balance, err
	}

	spended, err := wcs.withdrawalRepo.GetBalanceSpended(ctx, login)
	if err != nil {
		return balance, err
	}
	balance.Balance = added - spended
	balance.SpendedBonuses = spended

	return balance, nil
}

func (wcs *WithdrawalCreationService) GetWithdrawals(ctx context.Context, login string) ([]entities.FormatedWithdrawal, error) {
	withdrawals, err := wcs.withdrawalRepo.GetWithdrawals(ctx, login)
	if err != nil {
		return nil, err
	}

	fWithdrawals, err := wcs.formatDate(withdrawals)
	if err != nil {
		return nil, err
	}

	return fWithdrawals, nil
}

func (ocs *WithdrawalCreationService) formatDate(withdrawals []entities.Withdrawal) ([]entities.FormatedWithdrawal, error) {
	var fWithdrawals []entities.FormatedWithdrawal

	for _, w := range withdrawals {
		var fWithdrawal entities.FormatedWithdrawal
		fWithdrawal.OrderNum = w.OrderNum
		fWithdrawal.Bonuses = w.Bonuses
		fWithdrawal.FormatedDate = w.SpendDate.Format(time.RFC3339)

		fWithdrawals = append(fWithdrawals, fWithdrawal)
	}
	return fWithdrawals, nil
}
