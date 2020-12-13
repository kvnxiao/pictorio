package cookies

import (
	"net/http"
)

const (
	cookieUserID   = "uid"
	cookieUserName = "uname"
)

func SetUserID(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Name:  cookieUserID,
		Path:  "/",
		Value: encode(id),
	}
	http.SetCookie(w, cookie)
}

func GetUserID(r *http.Request) (string, error) {
	c, err := r.Cookie(cookieUserID)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", nil
		default:
			return "", err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func SetUserName(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:  cookieUserName,
		Path:  "/",
		Value: encode(name),
	}
	http.SetCookie(w, cookie)
}

func GetUserName(r *http.Request) (string, error) {
	c, err := r.Cookie(cookieUserName)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", nil
		default:
			return "", err
		}
	}
	value, err := decode(c.Value)
	if err != nil {
		return "", err
	}
	return string(value), nil
}
