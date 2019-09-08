package gonbayes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
)

var re = regexp.MustCompile("[^a-zA-Z 0-9]+")

func countWords(document string) (wordCount map[string]int) {
	document = clean(document)
	words := strings.Split(document, " ")
	wordCount = make(map[string]int)
	for _, word := range words {
		if !isStopWords(word) {
			key := stem(strings.ToLower(word))
			wordCount[key]++
		}
	}
	return
}

func clean(document string) string {
	return re.ReplaceAllString(strings.ToLower(document), "")
}

func isStopWords(word string) bool {
	_, ok := stopWords[word]
	return ok
}

func stem(word string) string {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil {
		// ignore error
		fmt.Printf("cannot stem word: %s\n", word)
		return word
	}
	return stemmed
}
