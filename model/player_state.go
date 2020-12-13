package model

type PlayerState struct {
	User         User `json:"user"`
	Points       int  `json:"points"`
	Wins         int  `json:"wins"`
	IsSpectator  bool `json:"isSpectator"`
	IsConnected  bool `json:"isConnected"`
	IsReady      bool `json:"isReady"`
	IsRoomLeader bool `json:"isRoomLeader"`
}
