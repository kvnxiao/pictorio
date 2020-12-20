package words

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
)

var wordBank []string

func init() {
	f, err := os.Open("assets/words.txt")
	if err != nil {
		log.Fatalln("Failed to read list of game words")
		return
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		wordBank = append(wordBank, strings.TrimSpace(scanner.Text()))
	}

	_ = f.Close()
}

func GenerateWord() string {
	return wordBank[rand.Intn(len(wordBank))]
}
