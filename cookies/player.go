package cookies

import (
	"net/http"
)

const (
	cookiePlayerID   = "pid"
	cookiePlayerName = "pname"
)

func SetPlayerID(w http.ResponseWriter, id string) {
	cookie := &http.Cookie{
		Name:  cookiePlayerID,
		Path:  "/",
		Value: encode(id),
	}
	http.SetCookie(w, cookie)
}

func GetPlayerID(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie(cookiePlayerID)
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

func SetPlayerName(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:  cookiePlayerName,
		Path:  "/",
		Value: encode(name),
	}
	http.SetCookie(w, cookie)
}

func GetPlayerName(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie(cookiePlayerName)
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
