package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vdgalyns/link-shortener/internal/cookies"
	"github.com/vdgalyns/link-shortener/internal/entities"
)

func ReadAndWriteCookieUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Попробовать прочитать значение куки
		_, errReadCookie := cookies.ReadSigned(r, "user_id")
		if errReadCookie != nil {
			if errors.Is(errReadCookie, http.ErrNoCookie) || errors.Is(errReadCookie, cookies.ErrInvalidValue) {
				// Создать уникальный идентификатор пользователя
				userID, errCreateUserID := entities.CreateUserID()
				if errCreateUserID != nil {
					http.Error(w, errCreateUserID.Error(), http.StatusBadRequest)
					return
				}
				newCookie := http.Cookie{
					Name:  "user_id",
					Value: userID,
					Path:  r.URL.Path,
				}
				// Записываем в куки результат
				errWriteCookie := cookies.WriteSigned(w, newCookie)
				if errWriteCookie != nil {
					http.Error(w, errWriteCookie.Error(), http.StatusBadRequest)
					return
				}
			} else {
				// Другая ошибка
				http.Error(w, errReadCookie.Error(), http.StatusBadRequest)
				return
			}
		}
		fmt.Println("OK", errReadCookie)
		// Отдаем в handler
		next.ServeHTTP(w, r)
	})
}
