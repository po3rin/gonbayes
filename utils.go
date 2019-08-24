package gonbayes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
)

func countWords(document string) (wordCount map[string]int) {
	replaced := regexp.MustCompile("[^a-zA-Z 0-9]+").ReplaceAllString(strings.ToLower(document), "")
	words := strings.Split(replaced, " ")
	wordCount = make(map[string]int)
	for _, word := range words {
		if _, ok := stopWords[word]; !ok {
			key := stem(strings.ToLower(word))
			wordCount[key]++
		}
	}
	return
}

func stem(word string) string {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil {
		// ignore error
		fmt.Println("Cannot stem word:", word)
		return word
	}
	return stemmed
}
