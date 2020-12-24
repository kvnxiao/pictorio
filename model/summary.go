package model

type WordSummary struct {
	Word           string   `json:"word"`
	WordLength     []int    `json:"wordLength"`
	WordSelections []string `json:"selections"`
}

type GameStateSummary struct {
	MaxRounds        int `json:"maxRounds"`
	MaxSelectionTime int `json:"maxSelectionTime"`
	MaxTurnTime      int `json:"maxTurnTime"`
	Round            int `json:"round"`
	TimeLeft         int `json:"timeLeft"`

	Status     GameStatus `json:"status"`
	TurnStatus TurnStatus `json:"turnStatus"`

	PlayerOrderIDs []string    `json:"playerOrderIds"`
	WordSummary    WordSummary `json:"words"`
}

type PlayersSummary struct {
	PlayerStates []PlayerState `json:"playerStates"`
	MaxPlayers   int           `json:"maxPlayers"`
}
