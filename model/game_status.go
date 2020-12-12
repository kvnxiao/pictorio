package model

type GameStatus int

const (
	StatusWaitingReadyUp GameStatus = iota
	StatusStarted
	StatusGameOver
)
