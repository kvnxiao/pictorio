package model

type GameStatus int

const (
	GameWaitingReadyUp GameStatus = iota
	GameStarted
	GameOver
)

type TurnStatus int

const (
	TurnSelection TurnStatus = iota
	TurnDrawing
	TurnEnded
)
