package model

type ChatEvent struct {
	User User `json:"user"`
	Message string `json:"message"`
}
