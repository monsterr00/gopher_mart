package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

/*
Получение списка загруженных номеров заказов
Хендлер: GET /api/user/orders.
Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым. Формат даты — RFC3339.
Доступные статусы обработки расчётов:
NEW — заказ загружен в систему, но не попал в обработку;
PROCESSING — вознаграждение за заказ рассчитывается;
INVALID — система расчёта вознаграждений отказала в расчёте;
PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
Формат запроса:
GET /api/user/orders HTTP/1.1
Content-Length: 0
Возможные коды ответа:
200 — успешная обработка запроса.
  Формат ответа:
  200 OK HTTP/1.1
  Content-Type: application/json
  ...

  [
      {
          "number": "9278923470",
          "status": "PROCESSED",
          "accrual": 500,
          "uploaded_at": "2020-12-10T15:15:45+03:00"
      },
      {
          "number": "12345678903",
          "status": "PROCESSING",
          "uploaded_at": "2020-12-10T15:12:01+03:00"
      },
      {
          "number": "346436439",
          "status": "INVALID",
          "uploaded_at": "2020-12-09T16:09:53+03:00"
      }
  ]

204 — нет данных для ответа.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/

func (h *OrderHandlers) Fetch(res http.ResponseWriter, req *http.Request) {
	login, _, _ := req.BasicAuth()
	ctx := req.Context()

	orders, err := h.oCrRepo.GetOrdersLogin(ctx, login)
	if errors.Is(err, i.ErrNoData) {
		http.Error(res, "No orders for user", http.StatusNoContent)
		return
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
