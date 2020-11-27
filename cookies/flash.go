package cookies

import (
	"net/http"
	"time"

	"github.com/kvnxiao/pictorio/model"
)

const (
	cookieFlashError = "errormsg"
)

func FlashError(w http.ResponseWriter, message string) {
	cookie := &http.Cookie{
		Name:  cookieFlashError,
		Path:  "/",
		Value: encode(message),
	}
	http.SetCookie(w, cookie)
}

func ReadError(w http.ResponseWriter, r *http.Request) (model.FlashMessage, error) {
	c, err := r.Cookie(cookieFlashError)
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
	deleteCookie := &http.Cookie{Name: cookieFlashError, MaxAge: -1, Path: "/", Expires: time.Unix(1, 0)}
	http.SetCookie(w, deleteCookie)
	return model.FlashMessage{Message: string(value), Type: model.FlashError}, nil
}
