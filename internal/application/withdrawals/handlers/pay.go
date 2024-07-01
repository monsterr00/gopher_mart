package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

/*
Запрос на списание средств
Хендлер: POST /api/user/balance/withdraw
Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя, в счёт оплаты которого списываются баллы.
Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.
Формат запроса:
POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

{
    "order": "2377225624",
    "sum": 751
}
Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.
Возможные коды ответа:
200 — успешная обработка запроса;
401 — пользователь не авторизован;
402 — на счету недостаточно средств;
422 — неверный номер заказа;
500 — внутренняя ошибка сервера.
*/

func (h *WithdrawalHandler) Withdraw(res http.ResponseWriter, req *http.Request) {
	var err error

	ctype := req.Header.Get("Content-Type")
	if ctype != "application/json" {
		http.Error(res, "Wrong content type", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var withdrawal entities.Withdrawal
	if err = json.Unmarshal(buf.Bytes(), &withdrawal); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var order entities.Order
	order.Login, _, _ = req.BasicAuth()
	order.OrderNum = withdrawal.OrderNum
	order.SpendedBonuses = withdrawal.Bonuses
	ctx := req.Context()

	isOk := order.ValidateOrderNum()
	if !isOk {
		http.Error(res, "Wrong order number format", http.StatusUnprocessableEntity)
		return
	}

	err = h.wCrRepo.SaveOrder(ctx, order)

	if errors.Is(err, i.ErrInsufFunds) {
		http.Error(res, err.Error(), http.StatusPaymentRequired)
		return
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("Order created/updated: %s", order.OrderNum)))
}
