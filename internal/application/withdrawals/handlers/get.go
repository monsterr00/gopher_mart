package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

/*
Получение текущего баланса пользователя
Хендлер: GET /api/user/balance.
Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.
Формат запроса:
GET /api/user/balance HTTP/1.1
Content-Length: 0
Возможные коды ответа:
200 — успешная обработка запроса.

	Формат ответа:
	200 OK HTTP/1.1
	Content-Type: application/json
	...

	{
	    "current": 500.5,
	    "withdrawn": 42
	}

401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func (h *WithdrawalHandler) GetBalance(res http.ResponseWriter, req *http.Request) {
	login, _, _ := req.BasicAuth()
	ctx := req.Context()

	balance, err := h.wCrRepo.GetBalance(ctx, login)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(balance)
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

/*
Получение информации о выводе средств
Хендлер: GET /api/user/withdrawals.
Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым. Формат даты — RFC3339.
Формат запроса:
GET /api/user/withdrawals HTTP/1.1
Content-Length: 0
Возможные коды ответа:
200 — успешная обработка запроса.
  Формат ответа:
  200 OK HTTP/1.1
  Content-Type: application/json
  ...

  [
      {
          "order": "2377225624",
          "sum": 500,
          "processed_at": "2020-12-09T16:09:57+03:00"
      }
  ]

204 — нет ни одного списания.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/

func (h *WithdrawalHandler) Fetch(res http.ResponseWriter, req *http.Request) {
	login, _, _ := req.BasicAuth()
	ctx := req.Context()

	withdrawals, err := h.wCrRepo.GetWithdrawals(ctx, login)
	if errors.Is(err, i.ErrNoData) {
		http.Error(res, "No withdrawals for user", http.StatusNoContent)
		return
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(withdrawals)
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
