package model

type GameStateStatus int

const (
	StatusWaitingReadyUp GameStateStatus = iota
	StatusStarted
	StatusGameOver
)
