package orders

import (
	"context"
	"time"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

func (ocs *OrderCreationService) GetOrdersLogin(ctx context.Context, login string) ([]entities.FormatedOrder, error) {
	orders, err := ocs.orderRepo.GetOrdersLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	fOrders, err := ocs.formatDate(orders)
	if err != nil {
		return nil, err
	}

	return fOrders, nil
}

func (ocs *OrderCreationService) formatDate(orders []entities.Order) ([]entities.FormatedOrder, error) {
	var fOrders []entities.FormatedOrder

	for _, o := range orders {
		var fOrder entities.FormatedOrder
		fOrder.OrderNum = o.OrderNum
		fOrder.Status = o.Status
		fOrder.AddedBonuses = o.AddedBonuses
		fOrder.CreatedAt = o.CreatedAt.Format(time.RFC3339)

		fOrders = append(fOrders, fOrder)
	}
	return fOrders, nil
}
