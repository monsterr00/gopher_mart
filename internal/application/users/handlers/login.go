package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/monsterr00/gopher_mart/internal/domain/entities"
)

/*
Аутентификация пользователя
Хендлер: POST /api/user/login.
Аутентификация производится по паре логин/пароль.
Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.
Формат запроса:
POST /api/user/login HTTP/1.1
Content-Type: application/json
...

{
    "login": "<login>",
    "password": "<password>"
}
Возможные коды ответа:
200 — пользователь успешно аутентифицирован;
400 — неверный формат запроса;
401 — неверная пара логин/пароль;
500 — внутренняя ошибка сервера.
*/

func (h *UserHandler) Login(res http.ResponseWriter, req *http.Request) {
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

	ctx := req.Context()
	err = h.uCrRepo.CheckAuth(ctx, user)
	if err != nil {
		res.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	auth := user.Login + ":" + user.Password
	res.Header().Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	res.WriteHeader(http.StatusOK)
}
