package model

type RoomResponse struct {
	RoomID string `json:"roomID"`
	Exists bool   `json:"exists"`
}
