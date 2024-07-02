package users

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (ucs *UserCreationService) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if ok {
			passwordHash, err := ucs.userRepo.GetPassword(r.Context(), username)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
			if err == nil {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
