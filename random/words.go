package random

import (
	"bufio"
	"math/rand"
	"os"
)

var words []string

func init() {
	f, err := os.Open("words.txt")
	if err != nil {
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
