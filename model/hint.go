package model

type Hint struct {
	Char      rune `json:"char"`
	WordIndex int  `json:"wordIndex"`
	CharIndex int  `json:"charIndex"`
}
