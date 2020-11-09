package cookies

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/kvnxiao/pictorio/model"
)

const (
	flashError = "errormsg"
)

func encode(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

func decode(input string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(input)
}

func FlashError(w http.ResponseWriter, message string) {
	cookie := &http.Cookie{
		Name:  flashError,
		Path:  "/",
		Value: encode(message),
	}
	http.SetCookie(w, cookie)
}

func ReadError(w http.ResponseWriter, r *http.Request) (model.FlashMessage, error) {
	c, err := r.Cookie(flashError)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return model.FlashMessage{}, nil
		default:
			return model.FlashMessage{}, err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return model.FlashMessage{}, err
	}
	deleteCookie := &http.Cookie{Name: flashError, MaxAge: -1, Path: "/", Expires: time.Unix(1, 0)}
	http.SetCookie(w, deleteCookie)
	return model.FlashMessage{Message: string(value), Type: model.FlashError}, nil
}
