package model

type RoomRequest struct {
	RoomID string `json:"roomID"`
}

type RoomResponse struct {
	RoomID string `json:"roomID"`
	Exists bool   `json:"exists"`
}
