package words

import (
	"math/rand"
	"strings"

	"github.com/kvnxiao/pictorio/model"
)

const censoredChar = "*"

var vowelsMap = map[rune]bool{
	'a': true,
	'e': true,
	'i': true,
	'o': true,
	'u': true,
	'y': true,
}

type GameWord struct {
	word       string
	wordLength []int
	hints      []model.Hint
}

func (w GameWord) Word() string {
	return w.word
}

func (w GameWord) WordLength() []int {
	return w.wordLength
}

func (w GameWord) Hints() []model.Hint {
	return w.hints
}

func (w GameWord) Censored() string {
	var censored []string
	for _, length := range w.wordLength {
		censored = append(censored, strings.Repeat(censoredChar, length))
	}
	return strings.Join(censored, " ")
}

func generateWordLength(word string) ([]int, []string) {
	split := strings.Fields(word)

	wordLengths := make([]int, len(split))
	for i := 0; i < len(split); i++ {
		wordLengths[i] = len(split[i])
	}
	return wordLengths, split
}

func isVowel(char rune) bool {
	return vowelsMap[char]
}

func generateAllHints(splitWords []string) []model.Hint {
	var hints []model.Hint

	for i := 0; i < len(splitWords); i++ {
		word := splitWords[i]
		for j, c := range word {
			char := c
			if !isVowel(char) {
				hints = append(hints, model.Hint{
					Char:      char,
					WordIndex: i,
					CharIndex: j,
				})
			}
		}
	}

	rand.Shuffle(len(hints), func(i, j int) {
		hints[i], hints[j] = hints[j], hints[i]
	})

	return hints
}

func NewGameWord(word string) GameWord {
	processedWord := strings.ToLower(word)
	wordLength, splitWords := generateWordLength(processedWord)
	hints := generateAllHints(splitWords)

	return GameWord{
		word:       processedWord,
		wordLength: wordLength,
		hints:      hints,
	}
}

func Censor(length int) string {
	return strings.Repeat(censoredChar, length)
}
