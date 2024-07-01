package handlers

import "net/http"

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
	login, password, _ := req.BasicAuth()

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set("Authorization", login+":"+password)
	res.WriteHeader(http.StatusOK)
}
