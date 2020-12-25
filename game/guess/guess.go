package guess

// PlayerGuesses represents the state of guesses for a single drawing turn (the drawer, and the rest of the players).
// The first player to guess correctly is awarded the most points, and the drawer is awarded some points (slightly less
// than the first guesser), while the rest of the players will be awarded less points if guessing correctly after the
// first guess has been established
type PlayerGuesses struct {
	guesses          map[string]bool
	firstGuessExists bool
}

func NewPlayerGuesses() *PlayerGuesses {
	return &PlayerGuesses{
		guesses:          make(map[string]bool),
		firstGuessExists: false,
	}
}

func (g *PlayerGuesses) HasGuessed(playerID string) bool {
	return g.guesses[playerID]
}

func (g *PlayerGuesses) AddGuessed(playerID string) (guesserPoints int, drawerPoints int) {
	g.guesses[playerID] = true
	if !g.firstGuessExists {
		g.firstGuessExists = true
		return 3, 2
	}
	return 1, 0
}
