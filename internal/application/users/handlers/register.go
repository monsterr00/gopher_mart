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
Регистрация пользователя

Хендлер: POST /api/user/register.
Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
После успешной регистрации должна происходить автоматическая аутентификация пользователя.
Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.

Формат запроса:
POST /api/user/register HTTP/1.1
Content-Type: application/json
...

	{
	    "login": "<login>",
	    "password": "<password>"
	}

Возможные коды ответа:
200 — пользователь успешно зарегистрирован и аутентифицирован;
400 — неверный формат запроса;
409 — логин уже занят;
500 — внутренняя ошибка сервера.
*/

func (h *UserHandler) Register(res http.ResponseWriter, req *http.Request) {
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

	var user entities.User
	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = user.ValidateUser()
	if errors.Is(err, i.ErrValidate) {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	err = h.uCrRepo.SaveUser(ctx, user)

	if errors.Is(err, i.ErrAlreadyExists) {
		http.Error(res, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("Authorization", user.Login+":"+user.Password)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("User registered:%s", user.Login)))
}
