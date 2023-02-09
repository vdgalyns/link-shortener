package cookies

import (
	"errors"
	"net/http"

	"github.com/vdgalyns/link-shortener/internal/entities"
)

const CookieNameUserID = "user_id"

func ReadAndCreateCookieUserID(w http.ResponseWriter, r *http.Request) (string, bool, error) {
	value, err := ReadSigned(r, CookieNameUserID)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) || errors.Is(err, ErrInvalidValue) {
			value, err = entities.CreateUserID()
			if err != nil {
				return value, false, err
			}
			err = WriteSigned(w, http.Cookie{Name: CookieNameUserID, Value: value})
			if err != nil {
				return value, false, err
			}
		}
		return value, false, nil
	}
	return value, true, nil
}
