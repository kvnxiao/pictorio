package model

type GameStatus int

const (
	StatusWaitingReadyUp GameStatus = iota
	StatusStarted
	StatusGameOver
)

type TurnStatus int

const (
	TurnSelection TurnStatus = iota
	TurnDrawing
	TurnEnded
)
