package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

func (ocs *OrderCreationService) SaveOrder(ctx context.Context, order entities.Order) error {
	err := ocs.orderRepo.CheckOrderDuplicate(ctx, order)
	if err != nil {
		return err
	}

	order.ID, _ = uuid.NewRandom()
	err = ocs.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		return err
	}

	err = ocs.registerOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

/*
Взаимодействие с системой расчёта начислений баллов лояльности
Для взаимодействия с системой доступен один хендлер:
GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.
Формат запроса:
GET /api/orders/{number} HTTP/1.1
Content-Length: 0
Возможные коды ответа:
200 — успешная обработка запроса.
  Формат ответа:
  200 OK HTTP/1.1
  Content-Type: application/json
  ...

  {
      "order": "<number>",
      "status": "PROCESSED",
      "accrual": 500
  }

  Поля объекта ответа:
order — номер заказа;
status — статус расчёта начисления:
REGISTERED — заказ зарегистрирован, но вознаграждение не рассчитано;
INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
PROCESSING — расчёт начисления в процессе;
PROCESSED — расчёт начисления окончен;
accrual — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
204 — заказ не зарегистрирован в системе расчёта.
429 — превышено количество запросов к сервису.
  Формат ответа:
  429 Too Many Requests HTTP/1.1
  Content-Type: text/plain
  Retry-After: 60

  No more than N requests per minute allowed

500 — внутренняя ошибка сервера.
Заказ может быть взят в расчёт в любой момент после его совершения. Время выполнения расчёта системой не регламентировано. Статусы INVALID и PROCESSED являются окончательными.
Общее количество запросов информации о начислении не ограничено.
*/

func (ocs *OrderCreationService) registerOrder(ctx context.Context, order entities.Order) error {
	req := ocs.client.R().
		SetHeader("Content-Type", "text/plain")

	requestURL := fmt.Sprintf("%s%s%s", ocs.host, "/api/orders/", order.OrderNum)

	res, err := req.Get(requestURL)
	if err != nil {
		return err
	}

	var regOrder entities.RegistredOrder
	if res.StatusCode() == http.StatusOK {
		if err = json.Unmarshal(res.Body(), &regOrder); err != nil {
			return err
		}
	}

	err = ocs.orderRepo.UpdateStatus(ctx, regOrder)
	if err != nil {
		return err
	}

	if regOrder.AddedBonuses != nil {
		err = ocs.orderRepo.UpdateBonuses(ctx, regOrder)
		if err != nil {
			return err
		}
	}
	return nil
}
