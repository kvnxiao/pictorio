package model

const FlashError = "error"

type FlashMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}
