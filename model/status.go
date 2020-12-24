package model

type GameStatus int

const (
	GameWaitingReadyUp GameStatus = iota
	GameStarted
	GameOver
)

type TurnStatus int

const (
	TurnNextPlayer TurnStatus = iota
	TurnSelection
	TurnDrawing
	TurnEnded
)
