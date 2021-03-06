// Package gonbayes is Simple Naive Bayes Classifier in Go.
package gonbayes

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"github.com/pkg/errors"
)

// Classifier is documents categories clasifier.
type Classifier struct {
	Words                  map[string]map[string]uint64
	TotalWords             uint64
	TotalDocsInCategories  map[string]uint64
	TotalDocs              uint64
	TotalWordsInCategories map[string]uint64
}

// NewClassifier inits classifier.
func NewClassifier(categories []string) *Classifier {
	c := &Classifier{
		Words:                  make(map[string]map[string]uint64),
		TotalDocsInCategories:  make(map[string]uint64),
		TotalWordsInCategories: make(map[string]uint64),
	}

	for _, category := range categories {
		c.Words[category] = make(map[string]uint64)
	}
	return c
}

// Train trains documents classifier.
func (c *Classifier) Train(category string, document string) {
	for word, count := range countWords(document) {
		c.Words[category][word] += uint64(count)
		c.TotalWordsInCategories[category] += uint64(count)
		c.TotalWords += uint64(count)
	}
	c.TotalDocsInCategories[category]++
	c.TotalDocs++
}

func (c *Classifier) pCategory(category string) float64 {
	// return float64(c.TotalDocsInCategories[category]) / float64(c.TotalDocs)

	// Take measures against underflow
	return math.Log(float64(c.TotalDocsInCategories[category]) / float64(c.TotalDocs))
}

func (c *Classifier) pDocCategory(category string, document string) float64 {
	// p := 1.0
	// for word := range countWords(document) {
	// 	p *= c.pWordCategory(category, word)
	// }
	// return p

	// Take measures against underflow
	var p float64
	for word := range countWords(document) {
		p += math.Log(c.pWordCategory(category, word))
	}
	return p
}

func (c *Classifier) pWordCategory(category string, word string) float64 {
	// return float64(c.Words[category][stem(word)]) / float64(c.TotalWordsInCategories[category])

	// Additive smoothings
	n := float64(c.Words[category][stem(word)] + 1)
	d := float64(c.TotalWordsInCategories[category] + c.TotalWords)
	return n / d
}

func (c *Classifier) pCategoryDocument(category string, document string) float64 {
	// return c.pDocCategory(category, document) * c.pCategory(category)

	// Take measures against underflow
	return c.pDocCategory(category, document) + c.pCategory(category)
}

// P is Probabilities of each categories.
func (c *Classifier) P(document string) map[string]float64 {
	p := make(map[string]float64)
	for category := range c.Words {
		p[category] = c.pCategoryDocument(category, document)
	}
	return p
}

// Classify classify documents.
func (c *Classifier) Classify(document string) string {
	prob := c.P(document)

	type sorted struct {
		category    string
		probability float64
	}

	var sp []sorted
	for c, p := range prob {
		sp = append(sp, sorted{c, p})
	}
	sort.Slice(sp, func(i, j int) bool {
		return sp[i].probability > sp[j].probability
	})

	return sp[0].category
}

// Encode trained Classifier
func (c *Classifier) Encode(fileName string) error {
	if c.TotalDocs == 0 {
		return errors.New("gonbayes: classifier is not trained yet")
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = gob.NewEncoder(f).Encode(&c)
	if err != nil {
		return fmt.Errorf("gonbayes: failed to encode classifier: %w", err)
	}
	return nil
}

// Decode CBOW output file to struct.
func (c *Classifier) Decode(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = gob.NewDecoder(f).Decode(c)
	if err != nil {
		return fmt.Errorf("gonbayes: failed to dencode file: %w", err)
	}
	return nil
}
