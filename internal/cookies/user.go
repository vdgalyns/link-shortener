package cookies

import (
	"errors"
	"net/http"

	"github.com/vdgalyns/link-shortener/internal/entities"
)

const CookieNameUserId = "user_id"

func ReadAndCreateCookieUserId(w http.ResponseWriter, r *http.Request) (string, bool, error) {
	value, err := ReadSigned(r, CookieNameUserId)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) || errors.Is(err, ErrInvalidValue) {
			value, err = entities.CreateUserID()
			if err != nil {
				return value, false, err
			}
			err = WriteSigned(w, http.Cookie{Name: CookieNameUserId, Value: value})
			if err != nil {
				return value, false, err
			}
			return value, false, nil
		} else {
			return value, false, err
		}
	}
	return value, true, nil
}
