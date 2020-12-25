package model

import (
	"github.com/kvnxiao/pictorio/game/settings"
)

type WordSummary struct {
	Word           string   `json:"word"`
	WordLength     []int    `json:"wordLength"`
	WordSelections []string `json:"selections"`
}

type GameStateSummary struct {
	Settings settings.GameSettings `json:"settings"`
	Round    int                   `json:"round"`
	TimeLeft int                   `json:"timeLeft"`

	Status     GameStatus `json:"status"`
	TurnStatus TurnStatus `json:"turnStatus"`

	PlayerOrderIDs []string    `json:"playerOrderIds"`
	WordSummary    WordSummary `json:"words"`
	Winners        []Winner    `json:"winners"`
}

type PlayersSummary struct {
	PlayerStates []PlayerState `json:"playerStates"`
	MaxPlayers   int           `json:"maxPlayers"`
}
