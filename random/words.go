package random

import (
	"bufio"
	"log"
	"math/rand"
	"os"
)

var words []string

func init() {
	f, err := os.Open("assets/words.txt")
	if err != nil {
		log.Fatalln("Failed to read list of game words")
		return
	}

	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	_ = f.Close()
}

func GenerateWord() string {
	return words[rand.Intn(len(words))]
}
