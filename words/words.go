package words

import (
	"strings"
)

const censoredChar = "*"

type GameWord struct {
	word       string
	wordLength []int
}

func (w GameWord) Word() string {
	return w.word
}

func (w GameWord) WordLength() []int {
	return w.wordLength
}

func (w GameWord) Censored() string {
	var censored []string
	for _, length := range w.wordLength {
		censored = append(censored, strings.Repeat(censoredChar, length))
	}
	return strings.Join(censored, " ")
}

func generateWordLength(word string) []int {
	split := strings.Fields(word)

	wordLengths := make([]int, len(split))
	for i := 0; i < len(split); i++ {
		wordLengths[i] = len(split[i])
	}
	return wordLengths
}

func NewGameWord(word string) GameWord {
	processedWord := strings.ToLower(word)
	return GameWord{
		word:       processedWord,
		wordLength: generateWordLength(processedWord),
	}
}

func Censor(length int) string {
	return strings.Repeat(censoredChar, length)
}
