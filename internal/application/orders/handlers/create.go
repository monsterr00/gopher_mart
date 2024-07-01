package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
	i "github.com/monsterr00/gopher_mart/internal/infrastructure"
)

/*
Загрузка номера заказа
Хендлер: POST /api/user/orders.
Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.
Формат запроса:
POST /api/user/orders HTTP/1.1
Content-Type: text/plain
...

12345678903
Возможные коды ответа:
200 — номер заказа уже был загружен этим пользователем;
202 — новый номер заказа принят в обработку;
400 — неверный формат запроса;
401 — пользователь не аутентифицирован;
409 — номер заказа уже был загружен другим пользователем;
422 — неверный формат номера заказа;
500 — внутренняя ошибка сервера.
*/

func (h *OrderHandlers) Create(res http.ResponseWriter, req *http.Request) {
	var err error

	ctype := req.Header.Get("Content-Type")
	if ctype != "text/plain" {
		http.Error(res, "Wrong content type", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var order entities.Order
	order.OrderNum = buf.String()
	order.Login, _, _ = req.BasicAuth()
	ctx := req.Context()

	isOk := order.ValidateOrderNum()
	if !isOk {
		http.Error(res, "Wrong order number format", http.StatusUnprocessableEntity)
		return
	}

	err = h.oCrRepo.SaveOrder(ctx, order)

	if errors.Is(err, i.ErrDuplicateUserOrder) {
		http.Error(res, err.Error(), http.StatusOK)
		return
	}

	if errors.Is(err, i.ErrDuplicateOrder) {
		http.Error(res, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusAccepted)
	res.Write([]byte(fmt.Sprintf("Order created: %s", order.OrderNum)))
}
