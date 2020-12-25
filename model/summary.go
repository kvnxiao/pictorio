package model

type WordSummary struct {
	Word           string   `json:"word"`
	WordLength     []int    `json:"wordLength"`
	WordSelections []string `json:"selections"`
}

type GameStateSummary struct {
	MaxRounds        int `json:"maxRounds"`
	MaxNextUpTime    int `json:"maxNextUpTime"`
	MaxSelectionTime int `json:"maxSelectionTime"`
	MaxDrawingTime   int `json:"maxDrawingTime"`
	MaxEndTime       int `json:"maxEndTime"`
	Round            int `json:"round"`
	TimeLeft         int `json:"timeLeft"`

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
