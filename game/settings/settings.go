package settings

const (
	MaxPlayers                   int = 8
	MaxRounds                    int = 2
	MaxSelectableWords           int = 3
	MaxTurnNextPlayerTimeSeconds int = 3
	MaxTurnSelectionTimeSeconds  int = 5
	MaxTurnDrawingTimeSeconds    int = 60
	MaxTurnDrawingTimeCutSeconds int = 10
	MaxTurnEndTimeSeconds        int = 5
	FirstHintTimeLeftSeconds     int = 20
	SecondHintTimeLeftSeconds    int = 15
	ThirdHintTimeLeftSeconds     int = 10
)

type GameSettings struct {
	MaxPlayers                   int   `json:"maxPlayers"`
	MaxRounds                    int   `json:"maxRounds"`
	MaxSelectableWords           int   `json:"maxSelectableWords"`
	MaxTurnNextPlayerTimeSeconds int   `json:"maxTurnNextSec"`
	MaxTurnSelectionTimeSeconds  int   `json:"maxTurnSelectSec"`
	MaxTurnDrawingTimeSeconds    int   `json:"maxTurnDrawSec"`
	MaxTurnDrawingTimeCutSeconds int   // do not include in JSON
	MaxTurnEndTimeSeconds        int   `json:"maxTurnEndSec"`
	HintSettings                 []int `json:"hints"`
}

func DefaultSettings() GameSettings {
	return GameSettings{
		MaxPlayers:                   MaxPlayers,
		MaxRounds:                    MaxRounds,
		MaxSelectableWords:           MaxSelectableWords,
		MaxTurnNextPlayerTimeSeconds: MaxTurnNextPlayerTimeSeconds,
		MaxTurnSelectionTimeSeconds:  MaxTurnSelectionTimeSeconds,
		MaxTurnDrawingTimeSeconds:    MaxTurnDrawingTimeSeconds,
		MaxTurnDrawingTimeCutSeconds: MaxTurnDrawingTimeCutSeconds,
		MaxTurnEndTimeSeconds:        MaxTurnEndTimeSeconds,
		HintSettings: []int{
			FirstHintTimeLeftSeconds,
			SecondHintTimeLeftSeconds,
			ThirdHintTimeLeftSeconds,
		},
	}
}
